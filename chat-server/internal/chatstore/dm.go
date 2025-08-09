package chatstore

import (
	"sync"
	"time"
)

// DM represents a direct message between users
type DM struct {
    Sender    string    `json:"sender"`
    Recipient string    `json:"recipient"`
    Body      string    `json:"body"`
    Timestamp time.Time `json:"timestamp"`
    Read      bool      `json:"read"`
    ID        string    `json:"id,omitempty"` // Optional unique ID
}

// DMStore provides thread-safe storage of direct messages
type DMStore struct {
    messages map[string][]DM // Key is username, value is slice of messages
    mu       sync.RWMutex
}

// Private singleton instance
var (
    instance *DMStore
    once     sync.Once
)

// GetDMStore returns the singleton instance of DMStore
func GetDMStore() *DMStore {
    once.Do(func() {
        instance = &DMStore{
            messages: make(map[string][]DM),
        }
    })
    return instance
}

// AddMessage stores a new direct message
func (s *DMStore) AddMessage(dm DM) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // // Store for sender
    // s.messages[dm.Sender] = append(s.messages[dm.Sender], dm)
    
    // Store for recipient
    s.messages[dm.Recipient] = append(s.messages[dm.Recipient], dm)
}

// GetMessages retrieves all messages for a user
func (s *DMStore) GetMessages(username string) []DM {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if msgs, ok := s.messages[username]; ok {
        // Return a copy to prevent external modification
        result := make([]DM, len(msgs))
        copy(result, msgs)
        return result
    }
    
    return []DM{}
}

// GetConversation retrieves all messages between two users
func (s *DMStore) GetConversation(user1, user2 string) []DM {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    messages := s.messages[user1]
    var conversation []DM
    
    for _, dm := range messages {
        if (dm.Sender == user1 && dm.Recipient == user2) || 
           (dm.Sender == user2 && dm.Recipient == user1) {
            conversation = append(conversation, dm)
        }
    }
    
    return conversation
}

// MarkAsRead marks messages as read
func (s *DMStore) MarkAsRead(username, otherUser string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if msgs, ok := s.messages[username]; ok {
        for i := range msgs {
            if msgs[i].Sender == otherUser && !msgs[i].Read {
                msgs[i].Read = true
            }
        }
        s.messages[username] = msgs
    }
}

// GetUnreadCount returns the number of unread messages for a user
func (s *DMStore) GetUnreadCount(username string) int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    count := 0
    if msgs, ok := s.messages[username]; ok {
        for _, dm := range msgs {
            if dm.Recipient == username && !dm.Read {
                count++
            }
        }
    }
    
    return count
}