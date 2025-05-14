package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

type CLI struct {
	conn             net.Conn
	reader           *Reader
	serverAddr       string
	subscribed       bool
	subscribedTopics []string
	scanner          *readline.Instance
}

func NewCLI(conn net.Conn, addr string) *CLI {
	scanner, err := readline.NewEx(&readline.Config{
		Prompt:                 addr + "> ",
		HistoryFile:            "/tmp/redgo_history.txt",
		HistorySearchFold:      true,
		AutoComplete:           nil,
		DisableAutoSaveHistory: true,
	})

	if err != nil {
		fmt.Println("Error initializing readline:", err)
		return nil
	}

	return &CLI{
		conn:       conn,
		reader:     NewReader(conn),
		serverAddr: addr,
		subscribed: false,
		scanner:    scanner,
	}
}

func (cli *CLI) Run() {
	cli.handleInterrupt()

	defer cli.scanner.Close()

	fmt.Println("Simple Redgo CLI (type 'exit' to quit)")

	for {
		if !cli.subscribed {
			fmt.Print(cli.serverAddr, "> ")
		}

		line, err := cli.scanner.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Handle Ctrl+C gracefully
				if cli.subscribed {
					fmt.Println("\nUnsubscribing...")
					cli.unsubscribe()
				} else {
					fmt.Println("\nExiting...")
					os.Exit(0)
				}
			} else {
				break
			}
		}

		if line == "exit" {
			break
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		cmd := strings.ToUpper(args[0])

		if cmd == "SUBSCRIBE" {
			cli.startSubscription(args)
			continue
		}

		if cmd == "CLEAR" {
			readline.ClearScreen(cli.scanner)
			continue
		}

		cli.sendAndPrint(args)
	}
}

func (cli *CLI) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range c {
			if cli.subscribed {
				fmt.Println("\nUnsubscribing...")
				cli.unsubscribe()
			} else {
				fmt.Println("\nExiting...")
				os.Exit(0)
			}
		}
	}()
}

func (cli *CLI) unsubscribe() {
	args := append([]string{"UNSUBSCRIBE"}, cli.subscribedTopics...)
	_, _ = cli.conn.Write(EncodeCommandAsRespString(args))
	cli.subscribed = false
	cli.subscribedTopics = nil
}

func (cli *CLI) startSubscription(args []string) {
	if len(args) < 2 {
		fmt.Println("SUBSCRIBE requires at least one channel")
		return
	}

	cli.subscribed = true
	cli.subscribedTopics = args[1:] // save channel list
	_, err := cli.conn.Write(EncodeCommandAsRespString(args))
	if err != nil {
		fmt.Println("Failed to subscribe:", err)
		cli.subscribed = false
		return
	}

	fmt.Println("Subscribed to channels:", strings.Join(cli.subscribedTopics, ", "))

	// Begin listening for responses while subscribed
	for cli.subscribed {
		msg, err := cli.reader.ParseFromRespString()
		if err != nil {
			fmt.Println("Error during subscription:", err)
			break
		}
		msg.WriteToConsole()
	}
}

func (cli *CLI) sendAndPrint(args []string) {
	_, err := cli.conn.Write(EncodeCommandAsRespString(args))
	if err != nil {
		fmt.Println("Failed to send command:", err)
		return
	}

	resp, err := cli.reader.ParseFromRespString()
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	if resp.Type() != R_ERROR {
		cli.scanner.SaveHistory(strings.Join(args, " "))
	}

	resp.WriteToConsole()
}
