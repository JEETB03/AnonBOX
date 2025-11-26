package core

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// P2PManager handles the libp2p host and networking
type P2PManager struct {
	Host          host.Host
	Ctx           context.Context
	DHT           *dht.IpfsDHT
	Peers         map[peer.ID]peer.AddrInfo
	PeerMutex     sync.RWMutex
	MsgChan       chan Message // Channel to send received messages to UI/CLI
	EncryptionKey []byte       // 32-byte key for AES-256
}

type Message struct {
	Sender  string
	Content string
	IsFile  bool
	Data    []byte
}

const (
	DiscoveryServiceTag = "anonbox-pub"
	ProtocolID          = "/anonbox/1.0.0"
)

func NewP2PManager(ctx context.Context) *P2PManager {
	return &P2PManager{
		Ctx:     ctx,
		Peers:   make(map[peer.ID]peer.AddrInfo),
		MsgChan: make(chan Message, 100),
	}
}

// SetKey sets the encryption key derived from a password
func (p *P2PManager) SetKey(password string) {
	// TODO: Use SHA-256 to derive key. For now, just using a placeholder or assuming the caller handles it.
	// In a real app, we'd import crypto/sha256.
	// For now, let's just assume the password IS the key if 32 bytes, or pad it.
	// This is a placeholder. The SetEncryptionKey method is preferred.
}

func (p *P2PManager) SetEncryptionKey(key []byte) {
	p.EncryptionKey = key
}

func (p *P2PManager) Start() error {
	// Create a new libp2p Host
	// Listen on all interfaces, random port
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		return err
	}
	p.Host = h

	// Set stream handler
	h.SetStreamHandler(ProtocolID, p.handleStream)
	h.SetStreamHandler(FileProtocolID, p.HandleFileStream)

	// Setup mDNS discovery
	if err := p.setupDiscovery(); err != nil {
		return err
	}

	log.Printf("Host created. ID: %s, Addrs: %v", h.ID(), h.Addrs())
	return nil
}

func (p *P2PManager) setupDiscovery() error {
	s := mdns.NewMdnsService(p.Host, DiscoveryServiceTag, p)
	return s.Start()
}

// HandlePeerFound is called by mDNS when a peer is found
func (p *P2PManager) HandlePeerFound(pi peer.AddrInfo) {
	p.PeerMutex.Lock()
	defer p.PeerMutex.Unlock()

	if _, exists := p.Peers[pi.ID]; !exists && pi.ID != p.Host.ID() {
		log.Printf("Found peer: %s", pi.ID)
		p.Peers[pi.ID] = pi
		// Connect to the peer
		if err := p.Host.Connect(p.Ctx, pi); err != nil {
			log.Printf("Failed to connect to peer %s: %v", pi.ID, err)
		} else {
			log.Printf("Connected to peer: %s", pi.ID)
		}
	}
}

func (p *P2PManager) handleStream(s network.Stream) {
	log.Printf("New stream from %s", s.Conn().RemotePeer())
	go p.readLoop(s)
}

func (p *P2PManager) readLoop(s network.Stream) {
	defer s.Close()
	buf := make([]byte, 4096) // Buffer size
	for {
		n, err := s.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Stream read error: %v", err)
			}
			return
		}

		data := buf[:n]
		var content string

		// Decrypt if key is set
		if len(p.EncryptionKey) > 0 {
			decrypted, err := Decrypt(data, p.EncryptionKey)
			if err != nil {
				log.Printf("Decryption failed: %v", err)
				continue
			}
			content = string(decrypted)
		} else {
			content = string(data)
		}

		log.Printf("Received: %s", content)
		p.MsgChan <- Message{
			Sender:  s.Conn().RemotePeer().String(),
			Content: content,
		}
	}
}

func (p *P2PManager) SendMessage(peerID peer.ID, msg string) error {
	s, err := p.Host.NewStream(p.Ctx, peerID, ProtocolID)
	if err != nil {
		return err
	}
	defer s.Close()

	data := []byte(msg)

	// Encrypt if key is set
	if len(p.EncryptionKey) > 0 {
		encrypted, err := Encrypt(data, p.EncryptionKey)
		if err != nil {
			return err
		}
		data = encrypted
	}

	_, err = s.Write(data)
	return err
}
