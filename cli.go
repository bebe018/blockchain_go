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
				log.Printf("讀取輸入錯誤: %v", err)
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
					fmt.Println("監聽器準備退出...")
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
		return &CommandEvent{Command: cmd, Success: true, Message: "退出信號"}
	}

	// 根據命令創建 FlagSet 並定義參數
	switch cmd {
	case "createblockchain":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		address := fs.String("address", "", "The address to send genesis block reward to")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("解析參數錯誤: %v", err)}
		}

		if *address == "" {
			return &CommandEvent{Command: cmd, Success: false, Message: "錯誤：'address' 參數是必需的。"}
		}

		cli.createBlockchain(*address, nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ 區塊鏈創建成功！"}

	case "getbalance":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		address := fs.String("address", "", "The address to get balance for")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("解析參數錯誤: %v", err)}
		}

		if *address == "" {
			return &CommandEvent{Command: cmd, Success: false, Message: "錯誤：'address' 參數是必需的。"}
		}

		cli.getBalance(*address, nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ 餘額查詢完成。"}

	case "send":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		from := fs.String("from", "", "Source wallet address")
		to := fs.String("to", "", "Destination wallet address")
		amount := fs.Int("amount", 0, "Amount to send")
		mine := fs.Bool("mine", false, "Mine immediately on the same node")

		if err := fs.Parse(cmdArgs); err != nil {
			return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("解析參數錯誤: %v", err)}
		}

		if *from == "" || *to == "" || *amount <= 0 {
			return &CommandEvent{Command: cmd, Success: false, Message: "錯誤：'from', 'to' 和 'amount' 參數是必需的，且 amount 必須大於 0。"}
		}

		cli.send(*from, *to, *amount, nodeID, *mine)
		return &CommandEvent{Command: cmd, Success: true, Message: fmt.Sprintf("✅ 交易發送成功: %d 單位從 %s 到 %s (挖礦: %t)。", *amount, *from, *to, *mine)}

	case "createwallet":
		cli.createWallet(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ 錢包創建完成。"}

	case "listaddresses":
		cli.listAddresses(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ 地址列表已印出。"}

	case "printchain":
		cli.printChain(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ 區塊鏈已印出。"}

	case "reindexutxo":
		cli.reindexUTXO(nodeID)
		return &CommandEvent{Command: cmd, Success: true, Message: "✅ UTXO 索引重建完成。"}

	default:
		cli.printUsage()
		return &CommandEvent{Command: cmd, Success: false, Message: fmt.Sprintf("未知命令: %s", cmd)}
	}
}
