package concurrent_sharded_map

import (
	"sync"
)

type ConcurrentShardedInt64Map map[int]*ShardInt64

type ShardInt64 struct {
	items map[int64]interface{}
	lock  *sync.RWMutex
}

func NewInt64() ConcurrentShardedInt64Map {
	c := make(ConcurrentShardedInt64Map, 256)

	for i := 0; i < 256; i++ {
		c[i] = &ShardInt64{
			items: make(map[int64]interface{}, 2048),
			lock:  new(sync.RWMutex),
		}
	}

	return c
}

func (c ConcurrentShardedInt64Map) Get(key int64) (interface{}, bool) {
	shard := c.getShard(key)
	shard.lock.RLock()

	defer shard.lock.RUnlock()

	if value, ok := shard.items[key]; ok {
		return value, true
	}

	return nil, false
}

func (c ConcurrentShardedInt64Map) Set(key int64, data interface{}) {
	shard := c.getShard(key)
	shard.lock.Lock()

	defer shard.lock.Unlock()

	shard.items[key] = data
}

func (c ConcurrentShardedInt64Map) Delete(key int64) {
	shard := c.getShard(key)
	shard.lock.Lock()

	defer shard.lock.Unlock()

	if _, ok := shard.items[key]; ok {
		delete(shard.items, key)
	}
}

func (c ConcurrentShardedInt64Map) getShard(key int64) (shard *ShardInt64) {
	return c[int(key)%256]
}
