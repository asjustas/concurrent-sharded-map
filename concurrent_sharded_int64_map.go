package concurrent_sharded_map

import (
	"hash/adler32"
	"reflect"
	"sync"
	"unsafe"
)

type ConcurrentShardedMap map[int]*Shard

type Shard struct {
	items map[string]interface{}
	lock  *sync.RWMutex
}

func New() ConcurrentShardedMap {
	c := make(ConcurrentShardedMap, 256)

	for i := 0; i < 256; i++ {
		c[i] = &Shard{
			items: make(map[string]interface{}, 2048),
			lock:  new(sync.RWMutex),
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
	checksum := adler32.Checksum(c.unsafeGetBytes(key))

	return c[int(checksum)%256]
}

// https://stackoverflow.com/questions/59209493/how-to-use-unsafe-get-a-byte-slice-from-a-string-without-memory-copy
func (c ConcurrentShardedMap) unsafeGetBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}
