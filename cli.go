package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type CommandEvent struct {
	Command string
	Success bool
	Message string
}

type CLI struct {
}

func (cli *CLI) StartListener() <-chan *CommandEvent {
	eventCh := make(chan *CommandEvent)

	go func() {
		defer close(eventCh)

		reader := bufio.NewReader(os.Stdin)

		cli.printUsage()

		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error occurs when reading signal: %v", err)
				break
			}

			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}

			event := cli.processCommand(input)
			if event != nil {
				eventCh <- event

				if event.Command == "quit" || event.Command == "exit" {
					fmt.Println("Server listener is ready to quit...")
					break
				}
			}
		}
	}()

	return eventCh
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) processCommand(input string) *CommandEvent {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	cmd := args[0]
	cmdArgs := args[1:]
	nodeID := os.Getenv("NODE_ID")

	if nodeID == "" {
		return &CommandEvent{Command: cmd, Success: false, Message: "NODE_ID env. var is not set!"}
	}

	if cmd == "quit" || cmd == "exit" {
		return &CommandEvent{Command: cmd, Success: true, Message: "quit sugnal"}
	}

	switch cmd {
	case "createblockchain":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		address := fs.String("address", "", "The address to send genesis block reward to")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("Error occurs when parsing parameters: %v", err)}
		}

		if *address == "" {
			return &CommandEvent{Command: cmd, Success: false, Message: "error：'address' parameters are needed"}
		}

		cli.createBlockchain(*address, nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "Blockchain created successfully"}

	case "getbalance":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		address := fs.String("address", "", "The address to get balance for")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("Error occurs when parsing parameters: %v", err)}
		}

		if *address == "" {
			return &CommandEvent{Command: cmd, Success: false, Message: "error:'address' parameters are needed"}
		}

		cli.getBalance(*address, nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "Get balance successfully"}

	case "send":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		from := fs.String("from", "", "Source wallet address")
		to := fs.String("to", "", "Destination wallet address")
		amount := fs.Int("amount", 0, "Amount to send")
		mine := fs.Bool("mine", false, "Mine immediately on the same node")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("Error occurs when parsing parameters: %v", err)}
		}

		if *from == "" || *to == "" || *amount <= 0 {
			return &CommandEvent{Command: cmd, Success: false, Message: "error:'from', 'to' and 'amount' parameters are needed and amout must larger than 0"}
		}

		cli.send(*from, *to, *amount, nodeID, *mine)
		return &CommandEvent{Command: cmd, Success: true, Message: fmt.Sprintf("Transaction sent successfully: %d units from %s to %s (Mining: %t).", *amount, *from, *to, *mine)}

	case "createwallet":
		cli.createWallet(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "Wallet created successfully"}

	case "listaddresses":
		cli.listAddresses(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "List of address has been printed successfully"}

	case "printchain":
		cli.printChain(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "Blockchain has been printed successfully"}

	case "reindexutxo":
		cli.reindexUTXO(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "UTXO reindex successfully"}

	default:
		cli.printUsage()
		return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("Unknown conmmand %s", cmd)}
	}
}
