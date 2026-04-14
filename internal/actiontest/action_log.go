package actiontest

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type ActionLogStore struct {
	mu   sync.RWMutex
	logs []ActionLog
}

var globalLogStore = &ActionLogStore{
	logs: make([]ActionLog, 0),
}

func AddLog(log ActionLog) {
	globalLogStore.mu.Lock()
	defer globalLogStore.mu.Unlock()

	// Generate a simple hex ID
	bytes := make([]byte, 4)
	rand.Read(bytes)
	log.ID = hex.EncodeToString(bytes)

	// Prepend to list
	globalLogStore.logs = append([]ActionLog{log}, globalLogStore.logs...)

	// Limit to max 100 logs
	if len(globalLogStore.logs) > 100 {
		globalLogStore.logs = globalLogStore.logs[:100]
	}
}

func GetLogs() []ActionLog {
	globalLogStore.mu.RLock()
	defer globalLogStore.mu.RUnlock()

	// Return a copy to avoid mutation outside
	logsCopy := make([]ActionLog, len(globalLogStore.logs))
	copy(logsCopy, globalLogStore.logs)
	return logsCopy
}
