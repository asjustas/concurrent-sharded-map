package concurrent_sharded_map

import "testing"

func BenchmarkGet(b *testing.B) {
	m := New()
	m.Set("test", "values")

	for n := 0; n < b.N; n++ {
		m.Get("test")
	}
}
