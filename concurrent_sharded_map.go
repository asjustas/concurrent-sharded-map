package concurrent_sharded_map

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

type ConcurrentShardedMap map[string]*Shard

type Shard struct {
	items map[string]interface{}
	lock *sync.RWMutex
}

func New() ConcurrentShardedMap {
	c := make(ConcurrentShardedMap, 256)

	for i := 0; i < 256; i++ {
		c[fmt.Sprintf("%02x", i)] = &Shard{
			items: make(map[string]interface{}, 2048),
			lock: new(sync.RWMutex),
		}
	}

	return c
}

func (c ConcurrentShardedMap) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)
	shard.lock.RLock()

	defer shard.lock.RUnlock()

	if value, ok := shard.items[key]; ok {
		return value, true
	}

	return nil, false
}

func (c ConcurrentShardedMap) Set(key string, data interface{}) {
	shard := c.getShard(key)
	shard.lock.Lock()

	defer shard.lock.Unlock()

	shard.items[key] = data
}

func (c ConcurrentShardedMap) Delete(key string) {
	shard := c.getShard(key)
	shard.lock.Lock()

	defer shard.lock.Unlock()

	if _, ok := shard.items[key]; ok {
		delete(shard.items, key)
	}
}

func (c ConcurrentShardedMap) getShard(key string) (shard *Shard) {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	shardKey :=  fmt.Sprintf("%x", hasher.Sum(nil))[0:2]

	return c[shardKey]
}
