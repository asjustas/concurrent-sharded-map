package concurrent_sharded_map

import (
	"testing"
)

func TestConcurrentShardedMap_Get(t *testing.T) {
	m := New()
	m.Set("test", "value")

	value, exists := m.Get("test")

	if value != "value" {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, "value", true)
	}
}

func TestConcurrentShardedMap_Get_NonExistingKey(t *testing.T) {
	m := New()
	m.Set("test", "value")

	value, exists := m.Get("non_existing")

	if exists {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, nil, false)
	}

	if value != nil {
		t.Errorf("Get() = %v, %v, want %v, %v", value, exists, nil, false)
	}
}
