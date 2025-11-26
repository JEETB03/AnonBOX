package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"AnonBOX/pkg/core"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

var (
	p2pMgr *core.P2PManager
	ctx    context.Context
)

func main() {
	ctx = context.Background()
	p2pMgr = core.NewP2PManager(ctx)

	var rootCmd = &cobra.Command{
		Use:   "piratecove",
		Short: "P2P Amnesic Chatroom",
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the P2P node",
		Run:   runStart,
	}

	startCmd.Flags().StringP("password", "p", "", "Password for encryption (Vault)")

	rootCmd.AddCommand(startCmd)
	rootCmd.Execute()
}

func runStart(cmd *cobra.Command, args []string) {
	password, _ := cmd.Flags().GetString("password")
	if password != "" {
		hash := sha256.Sum256([]byte(password))
		p2pMgr.SetEncryptionKey(hash[:])
		fmt.Println("🔒 Encryption enabled with provided password.")
	} else {
		fmt.Println("⚠️  No password provided. Messages will be unencrypted (or use transport security only).")
	}

	if err := p2pMgr.Start(); err != nil {
		fmt.Printf("Error starting P2P host: %v\n", err)
		return
	}

	fmt.Println("🏴‍☠️  PirateCove Node Started!")
	fmt.Printf("ID: %s\n", p2pMgr.Host.ID())
	fmt.Println("Type 'help' for commands.")

	// Handle incoming messages
	go func() {
		for msg := range p2pMgr.MsgChan {
			if msg.IsFile {
				fmt.Printf("\n📦 [FILE] %s: %s\n> ", msg.Sender, msg.Content)
			} else {
				fmt.Printf("\n💬 [%s]: %s\n> ", msg.Sender, msg.Content)
			}
		}
	}()

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "peers":
			p2pMgr.PeerMutex.RLock()
			fmt.Println("Connected Peers:")
			for id := range p2pMgr.Peers {
				fmt.Printf("- %s\n", id)
			}
			p2pMgr.PeerMutex.RUnlock()
		case "chat":
			if len(parts) < 3 {
				fmt.Println("Usage: chat <peerID> <message>")
				continue
			}
			targetID, err := peer.Decode(parts[1])
			if err != nil {
				fmt.Printf("Invalid peer ID: %v\n", err)
				continue
			}
			msg := strings.Join(parts[2:], " ")
			if err := p2pMgr.SendMessage(targetID, msg); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
			} else {
				fmt.Println("Sent.")
			}
		case "share":
			if len(parts) < 3 {
				fmt.Println("Usage: share <peerID> <filePath>")
				continue
			}
			targetID, err := peer.Decode(parts[1])
			if err != nil {
				fmt.Printf("Invalid peer ID: %v\n", err)
				continue
			}
			filePath := parts[2]
			if err := p2pMgr.SendFile(targetID, filePath); err != nil {
				fmt.Printf("Error sending file: %v\n", err)
			} else {
				fmt.Println("File sent.")
			}
		case "broadcast":
			// Simple broadcast to all known peers
			msg := strings.Join(parts[1:], " ")
			p2pMgr.PeerMutex.RLock()
			for id := range p2pMgr.Peers {
				go p2pMgr.SendMessage(id, msg)
			}
			p2pMgr.PeerMutex.RUnlock()
			fmt.Println("Broadcast sent.")
		case "exit", "quit":
			fmt.Println("Exiting...")
			return
		case "help":
			fmt.Println("Commands: peers, chat <id> <msg>, share <id> <file>, broadcast <msg>, exit")
		default:
			fmt.Println("Unknown command. Try 'help'.")
		}
		fmt.Print("> ")
	}
}
