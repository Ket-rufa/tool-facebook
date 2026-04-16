package fbdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store defines operations for Facebook Data Store
type Store interface {
	Get(uid string) (*FBAccountData, error)
	Save(uid string, data FBAccountData) error
	ListAllUIDs() ([]string, error)
	Remove(uid string) error
}

// JSONStore is a file-backed JSON data store mapping UIDs to FBAccountData
type JSONStore struct {
	mu       sync.RWMutex
	filePath string
	// map[UID] FBAccountData
	cache map[string]FBAccountData
}

// NewJSONStore creates or loads the facebook_data.json store
func NewJSONStore(dir string) (*JSONStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục data: %w", err)
	}
	
	filePath := filepath.Join(dir, "facebook_data.json")
	s := &JSONStore{
		filePath: filePath,
		cache:    make(map[string]FBAccountData),
	}
	
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("không thể đọc dữ liệu FB data: %w", err)
	}
	
	return s, nil
}

func (s *JSONStore) load() error {
	b, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.cache)
}

func (s *JSONStore) persist() error {
	b, err := json.MarshalIndent(s.cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, b, 0644)
}

// Get retrieves data for a specific UID
func (s *JSONStore) Get(uid string) (*FBAccountData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	val, ok := s.cache[uid]
	if !ok {
		return nil, nil
	}
	// Return a copy
	res := val
	return &res, nil
}

// Save inserts or updates data for a specific UID
func (s *JSONStore) Save(uid string, data FBAccountData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.cache[uid] = data
	return s.persist()
}

// ListAllUIDs returns a list of all UIDs currently in the store
func (s *JSONStore) ListAllUIDs() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	uids := make([]string, 0, len(s.cache))
	for k := range s.cache {
		uids = append(uids, k)
	}
	return uids, nil
}

// Remove deletes the data for a specific UID
func (s *JSONStore) Remove(uid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.cache[uid]; exists {
		delete(s.cache, uid)
		return s.persist()
	}
	return nil
}
