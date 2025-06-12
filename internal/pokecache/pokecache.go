package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry map[string]cacheEntry
	m     sync.Mutex
}

type cacheEntry struct {
	CreateAt time.Time
	Value    []byte
}
