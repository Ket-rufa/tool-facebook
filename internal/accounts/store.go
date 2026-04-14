package accounts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store là interface local persistence - dễ swap sang SQLite sau này
type Store interface {
	List() ([]AccountProfile, error)
	Add(profile AccountProfile) error
	Remove(id string) error
	SetDefault(id string) error
	GetDefault() (*AccountProfile, error)
	UpdateSessionStatus(id string, status SessionStatus, checkedAt string) error
	Get(id string) (*AccountProfile, error)
}

// JSONStore lưu dữ liệu trong file JSON với locking
type JSONStore struct {
	mu       sync.RWMutex
	filePath string
	cache    []AccountProfile
}

func NewJSONStore(dir string) (*JSONStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục data: %w", err)
	}
	filePath := filepath.Join(dir, "accounts.json")
	s := &JSONStore{
		filePath: filePath,
		cache:    []AccountProfile{},
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("không thể đọc dữ liệu tài khoản: %w", err)
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

func (s *JSONStore) List() ([]AccountProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]AccountProfile, len(s.cache))
	copy(result, s.cache)
	return result, nil
}

func (s *JSONStore) Add(profile AccountProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Chỉ có 1 account mặc định
	if len(s.cache) == 0 {
		profile.IsDefault = true
	} else {
		profile.IsDefault = false
	}
	// Prepend để hiển thị mới nhất lên đầu
	s.cache = append([]AccountProfile{profile}, s.cache...)
	return s.persist()
}

func (s *JSONStore) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, a := range s.cache {
		if a.ID == id {
			s.cache = append(s.cache[:i], s.cache[i+1:]...)
			return s.persist()
		}
	}
	return nil
}

func (s *JSONStore) SetDefault(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.cache {
		s.cache[i].IsDefault = (s.cache[i].ID == id)
	}
	return s.persist()
}

func (s *JSONStore) GetDefault() (*AccountProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.cache {
		if a.IsDefault {
			acc := a
			return &acc, nil
		}
	}
	return nil, nil
}

func (s *JSONStore) UpdateSessionStatus(id string, status SessionStatus, checkedAt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if checkedAt == "" {
		checkedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	for i, a := range s.cache {
		if a.ID == id {
			s.cache[i].SessionStatus = string(status)
			s.cache[i].LastCheckedAt = checkedAt
			return s.persist()
		}
	}
	return nil
}

func (s *JSONStore) Get(id string) (*AccountProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.cache {
		if a.ID == id {
			acc := a
			return &acc, nil
		}
	}
	return nil, nil
}
