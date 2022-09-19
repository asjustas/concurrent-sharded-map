package concurrent_sharded_map

import (
	"testing"
)

func TestConcurrentShardedInt64Map_Get(t *testing.T) {
	m := NewInt64()
	m.Set(1, "value")

	value, exists := m.Get(1)

	if value != "value" {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, "value", true)
	}
}

func TestConcurrentShardedInt64Map_Get_NonExistingKey(t *testing.T) {
	m := NewInt64()
	m.Set(1, "value")

	value, exists := m.Get(2)

	if exists {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, nil, false)
	}

	if value != nil {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, nil, false)
	}
}
