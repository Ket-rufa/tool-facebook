package actiontest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type actionConfig struct {
	ReactDocID string `json:"react_doc_id"`
}

type ConfigStore struct {
	mu       sync.RWMutex
	filePath string
	cache    actionConfig
}

func NewConfigStore(dataDir string) (*ConfigStore, error) {
	if strings.TrimSpace(dataDir) == "" {
		return nil, errors.New("data dir khong hop le")
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("khong tao duoc data dir: %w", err)
	}

	store := &ConfigStore{
		filePath: filepath.Join(dataDir, "action_settings.json"),
		cache:    actionConfig{},
	}
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("khong doc duoc action settings: %w", err)
	}
	return store, nil
}

func (s *ConfigStore) load() error {
	b, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.cache)
}

func (s *ConfigStore) persistLocked() error {
	b, err := json.MarshalIndent(s.cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, b, 0644)
}

func (s *ConfigStore) GetReactDocID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return strings.TrimSpace(s.cache.ReactDocID)
}

func (s *ConfigStore) SetReactDocID(docID string) error {
	docID = strings.TrimSpace(docID)
	if docID != "" && !isNumericString(docID) {
		return errors.New("doc_id chi duoc chua so")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache.ReactDocID = docID
	return s.persistLocked()
}

func isNumericString(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
