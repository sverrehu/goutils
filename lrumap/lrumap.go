package lrumap

import (
	"sync"
	"time"
)

type cacheEntry struct {
	value        interface{}
	whenLastUsed time.Time
	whenExpires  time.Time
}

type LRUMap struct {
	entries map[string]*cacheEntry
	maxSize int
	ttl     time.Duration
	mutex   sync.Mutex
}

func New(maxSize int, ttl time.Duration) (m *LRUMap) {
	return &LRUMap{entries: make(map[string]*cacheEntry, maxSize), maxSize: maxSize, ttl: ttl}
}

func (m *LRUMap) Put(key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.makeRoomForOne()
	v, ok := m.entries[key]
	if !ok {
		v = &cacheEntry{}
	}
	now := time.Now()
	v.value = value
	v.whenLastUsed = now
	v.whenExpires = now.Add(m.ttl)
	m.entries[key] = v
}

func (m *LRUMap) Get(key string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, ok := m.entries[key]; ok {
		now := time.Now()
		if now.After(v.whenExpires) {
			delete(m.entries, key)
			return nil
		}
		v.whenLastUsed = now
		return v.value
	}
	return nil
}

func (m *LRUMap) Remove(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.entries, key)
}

func (m *LRUMap) makeRoomForOne() {
	// map must be locked before calling
	if len(m.entries) < m.maxSize {
		return
	}
	// The ultimate goal is to delete the least recently used (LRU) entry.
	// But if an expired entry is discovered along the way,
	// delete this one instead and bail out early.
	var lruKey string
	var lruTime time.Time
	first := true
	now := time.Now()
	for k, v := range m.entries {
		if now.After(v.whenExpires) {
			delete(m.entries, k)
			return
		}
		if first || v.whenExpires.Before(lruTime) {
			lruKey = k
			lruTime = v.whenExpires
			first = false
		}
	}
	delete(m.entries, lruKey)
}
