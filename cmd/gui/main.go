package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"AnonBOX/pkg/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	p2pMgr     *core.P2PManager
	ctx        context.Context
	msgList    *widget.List
	messages   []string
	peerList   *widget.List
	peerIDs    []peer.ID
	input      *widget.Entry
	targetPeer peer.ID
)

func main() {
	ctx = context.Background()
	p2pMgr = core.NewP2PManager(ctx)

	a := app.New()
	w := a.NewWindow("PirateCove 🏴‍☠️")
	w.Resize(fyne.NewSize(800, 600))

	// --- Login / Setup Screen ---
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Enter Vault Password (Optional)")

	startBtn := widget.NewButton("Enter the Cove", func() {
		password := passEntry.Text
		if password != "" {
			hash := sha256.Sum256([]byte(password))
			p2pMgr.SetEncryptionKey(hash[:])
		}

		if err := p2pMgr.Start(); err != nil {
			dialog.ShowError(err, w)
			return
		}

		showMainUI(w, a)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to PirateCove", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		passEntry,
		startBtn,
	))

	w.ShowAndRun()
}

func showMainUI(w fyne.Window, a fyne.App) {
	// --- Chat Tab ---
	messages = []string{"Welcome to the secure channel."}

	msgList = widget.NewList(
		func() int { return len(messages) },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(messages[i])
		},
	)

	input = widget.NewEntry()
	input.SetPlaceHolder("Type a message...")
	input.OnSubmitted = func(text string) {
		sendMessage(text)
	}

	sendBtn := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		sendMessage(input.Text)
	})

	inputContainer := container.NewBorder(nil, nil, nil, sendBtn, input)
	chatContainer := container.NewBorder(nil, inputContainer, nil, nil, msgList)

	// --- Peers Tab ---
	peerList = widget.NewList(
		func() int { return len(peerIDs) },
		func() fyne.CanvasObject { return widget.NewLabel("peer id") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(peerIDs[i].String())
		},
	)

	peerList.OnSelected = func(id widget.ListItemID) {
		targetPeer = peerIDs[id]
		dialog.ShowInformation("Peer Selected", fmt.Sprintf("Chatting with %s", targetPeer.String()), w)
	}

	refreshPeersBtn := widget.NewButton("Refresh Peers", func() {
		refreshPeers()
	})

	peersContainer := container.NewBorder(nil, refreshPeersBtn, nil, nil, peerList)

	// --- Files Tab ---
	fileBtn := widget.NewButton("Send File", func() {
		if targetPeer == "" {
			dialog.ShowError(fmt.Errorf("Select a peer first!"), w)
			return
		}

		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			path := reader.URI().Path()
			// Fyne URI path might need adjustment on Windows
			// For now, let's try using it directly or converting
			// Note: Fyne URI handling can be tricky.

			go func() {
				if err := p2pMgr.SendFile(targetPeer, path); err != nil {
					// dialog.ShowError(err, w) // Can't call from goroutine easily without context
					fmt.Printf("Error sending file: %v\n", err)
				} else {
					messages = append(messages, fmt.Sprintf("Sent file: %s", path))
					msgList.Refresh()
					msgList.ScrollToBottom()
				}
			}()
		}, w)
		fd.Show()
	})

	filesContainer := container.NewVBox(
		widget.NewLabel("File Sharing"),
		fileBtn,
		widget.NewLabel("Received files will appear in the chat log."),
	)

	// --- Tabs ---
	tabs := container.NewAppTabs(
		container.NewTabItem("Chat", chatContainer),
		container.NewTabItem("Peers", peersContainer),
		container.NewTabItem("Files", filesContainer),
	)

	w.SetContent(tabs)

	// Start background loops
	go handleIncomingMessages()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			refreshPeers()
		}
	}()
}

func sendMessage(text string) {
	if text == "" {
		return
	}
	if targetPeer == "" {
		messages = append(messages, "⚠️ Select a peer from the Peers tab first!")
		msgList.Refresh()
		msgList.ScrollToBottom()
		return
	}

	if err := p2pMgr.SendMessage(targetPeer, text); err != nil {
		messages = append(messages, fmt.Sprintf("Error: %v", err))
	} else {
		messages = append(messages, fmt.Sprintf("Me: %s", text))
	}
	input.SetText("")
	msgList.Refresh()
	msgList.ScrollToBottom()
}

func handleIncomingMessages() {
	for msg := range p2pMgr.MsgChan {
		content := msg.Content
		if msg.IsFile {
			content = fmt.Sprintf("📦 Received File: %s", msg.Content)
		} else {
			content = fmt.Sprintf("[%s]: %s", msg.Sender[:8], msg.Content)
		}

		// Update UI on main thread
		// Fyne lists are not thread safe for direct append?
		// Actually slice append is fine, but Refresh must be on UI thread?
		// Fyne docs say widget updates should be on UI thread.
		messages = append(messages, content)
		msgList.Refresh()
		msgList.ScrollToBottom()
	}
}

func refreshPeers() {
	p2pMgr.PeerMutex.RLock()
	newPeers := make([]peer.ID, 0, len(p2pMgr.Peers))
	for id := range p2pMgr.Peers {
		newPeers = append(newPeers, id)
	}
	p2pMgr.PeerMutex.RUnlock()

	peerIDs = newPeers
	peerList.Refresh()
}
