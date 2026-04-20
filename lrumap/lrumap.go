package lrumap

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value        interface{}
	WhenLastUsed time.Time
	WhenExpires  time.Time
}

type LRUMap struct {
	Entries map[string]*CacheEntry
	MaxSize int
	TTL     time.Duration
	mutex   sync.Mutex
}

func New(maxSize int, ttl time.Duration) (m *LRUMap) {
	return &LRUMap{Entries: make(map[string]*CacheEntry, maxSize), MaxSize: maxSize, TTL: ttl}
}

func (m *LRUMap) Put(key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.makeRoomForOne()
	v, ok := m.Entries[key]
	if !ok {
		v = &CacheEntry{}
	}
	now := time.Now()
	v.Value = value
	v.WhenLastUsed = now
	v.WhenExpires = now.Add(m.TTL)
	m.Entries[key] = v
}

func (m *LRUMap) Get(key string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, ok := m.Entries[key]; ok {
		now := time.Now()
		if now.After(v.WhenExpires) {
			delete(m.Entries, key)
			return nil
		}
		v.WhenLastUsed = now
		return v.Value
	}
	return nil
}

func (m *LRUMap) Remove(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.Entries, key)
}

func (m *LRUMap) makeRoomForOne() {
	// map must be locked before calling
	if len(m.Entries) < m.MaxSize {
		return
	}
	// The ultimate goal is to delete the least recently used (LRU) entry.
	// But if an expired entry is discovered along the way,
	// delete this one instead and bail out early.
	var lruKey string
	var lruTime time.Time
	first := true
	now := time.Now()
	for k, v := range m.Entries {
		if now.After(v.WhenExpires) {
			delete(m.Entries, k)
			return
		}
		if first || v.WhenExpires.Before(lruTime) {
			lruKey = k
			lruTime = v.WhenExpires
			first = false
		}
	}
	delete(m.Entries, lruKey)
}
