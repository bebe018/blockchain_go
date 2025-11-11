package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"os"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// SerializableWallet Solve gob: type elliptic.p256Curve has no exported fields
type SerializableWallet struct {
	D         *big.Int // D means private key
	X, Y      *big.Int // X, Y mean the curve coordinate to the public key
	PublicKey []byte
}

// NewWallets creates Wallets and fills it from a file if it exists
func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile(nodeID)

	return &wallets, err
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a Wallet by its address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile(nodeID string) error {

	walletFile := fmt.Sprintf(walletFile, nodeID)
	//fmt.Println(walletFile)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile) // use os.ReadFile to replace ioutil.ReadFile since ioutil package has been deprecated
	if err != nil {
		log.Panic(err)
	}

	var wallets map[string]SerializableWallet
	//gob.Register(elliptic.P256())
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = make(map[string]*Wallet)
	//ws.Wallets = wallets.Wallets
	for k, v := range wallets {
		ws.Wallets[k] = &Wallet{
			PrivateKey: ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}

	return nil
}

// SaveToFile saves wallets to a file
func (ws Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)
	//fmt.Println(walletFile)

	gob.Register(SerializableWallet{})

	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.Wallets {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}
