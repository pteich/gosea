package infrastructure

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	errorNotFound    = errors.New("cache item not found")
	errorTTLExceeded = errors.New("cache item too old")
)

type cacheItem struct {
	data []byte
	ts   time.Time
}

type RequestCache struct {
	maxTTL       time.Duration
	cache        map[string]cacheItem
	protectCache sync.RWMutex
}

func NewRequestCache(ttl time.Duration) *RequestCache {
	return &RequestCache{
		maxTTL: ttl,
		cache:  make(map[string]cacheItem),
	}
}

func (rc *RequestCache) Inject(cfg *struct {
	DefaultTTL float64 `inject:"config:seabackend.defaultCacheTTL"`
}) {
	if cfg != nil {
		rc.maxTTL = time.Duration(cfg.DefaultTTL) * time.Second
	}
	rc.cache = make(map[string]cacheItem)
}

func (rc *RequestCache) Set(key string, data interface{}) error {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return err
	}

	rc.protectCache.Lock()
	defer rc.protectCache.Unlock()

	rc.cache[key] = cacheItem{
		ts:   time.Now(),
		data: buf.Bytes(),
	}

	return nil
}

func (rc *RequestCache) Get(key string, data interface{}) error {
	rc.protectCache.RLock()
	defer rc.protectCache.RUnlock()

	cacheItem, found := rc.cache[key]
	if !found {
		return errorNotFound
	}

	if time.Now().Sub(cacheItem.ts) > rc.maxTTL {
		return errorTTLExceeded
	}

	buf := bytes.NewBuffer(cacheItem.data)
	err := gob.NewDecoder(buf).Decode(data)
	if err != nil {
		return fmt.Errorf("could not decode data: %w", err)
	}

	return nil
}
