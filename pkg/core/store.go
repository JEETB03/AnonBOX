package core

import (
	"sync"
	"time"
)

type ChatMessage struct {
	Sender    string
	Content   string
	Timestamp time.Time
	IsMine    bool
}

// Store handles in-memory storage of messages (Amnesic)
type Store struct {
	Messages []ChatMessage
	mu       sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		Messages: make([]ChatMessage, 0),
	}
}

func (s *Store) AddMessage(msg ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, msg)
}

func (s *Store) GetMessages() []ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to avoid race conditions if the caller modifies it
	msgs := make([]ChatMessage, len(s.Messages))
	copy(msgs, s.Messages)
	return msgs
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = make([]ChatMessage, 0)
}
