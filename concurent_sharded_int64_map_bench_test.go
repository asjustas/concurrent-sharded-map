package concurrent_sharded_map

import "testing"

func BenchmarkConcurrentShardedInt64Map_GetGet(b *testing.B) {
	m := NewInt64()
	m.Set(1, "values")

	for n := 0; n < b.N; n++ {
		m.Get(1)
	}
}
