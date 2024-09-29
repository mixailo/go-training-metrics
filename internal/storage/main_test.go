package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var storage MetricsStorage

func init() {
	storage = NewStorage()
	storage.UpdateCounter("test_counter", 1)
	storage.UpdateGauge("test_gauge", 1)
}

func TestMemStorage_GetCounter(t *testing.T) {
	tests := []struct {
		name   string
		index  string
		value  int64
		exists bool
	}{
		{
			name:   "existing",
			index:  "test_counter",
			value:  1,
			exists: true,
		},
		{
			name:   "not existing",
			index:  "test_counter_2",
			value:  1,
			exists: false,
		},
	}
	for _, tt := range tests {
		value, exists := storage.GetCounter(tt.index)
		if tt.exists {
			assert.True(t, exists)
			assert.Equal(t, tt.value, value)
		} else {
			assert.False(t, exists)
		}
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name   string
		index  string
		value  float64
		exists bool
	}{
		{
			name:   "existing",
			index:  "test_gauge",
			value:  1,
			exists: true,
		},
		{
			name:   "not existing",
			index:  "test_gauge_2",
			value:  1,
			exists: false,
		},
	}
	for _, tt := range tests {
		value, exists := storage.GetGauge(tt.index)
		if tt.exists {
			assert.True(t, exists)
			assert.Equal(t, tt.value, value)
		} else {
			assert.False(t, exists)
		}
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	value, ok := storage.GetCounter("test_counter")
	if !ok {
		t.Error("counter does not exist")
	} else {
		storage.UpdateCounter("test_counter", 1)
		newValue, ok := storage.GetCounter("test_counter")
		if !ok {
			t.Error("cannot get value after update")
		}
		assert.Equal(t, newValue, value+1)

		storage.UpdateCounter("test_counter", -2)
		newValue, ok = storage.GetCounter("test_counter")
		if !ok {
			t.Error("cannot get value after update")
		}
		assert.Equal(t, newValue, value-1)
	}

}

func TestMemStorage_UpdateGauge(t *testing.T) {
	value, ok := storage.GetCounter("test_counter")
	if !ok {
		t.Error("counter does not exist")
	} else {
		storage.UpdateCounter("test_counter", 1)
		newValue, ok := storage.GetCounter("test_counter")
		if !ok {
			t.Error("cannot get value after update")
		}
		assert.Equal(t, newValue, value+1)

		storage.UpdateCounter("test_counter", -2)
		newValue, ok = storage.GetCounter("test_counter")
		if !ok {
			t.Error("cannot get value after update")
		}
		assert.Equal(t, newValue, value-1)
	}
}

func TestNewStorage(t *testing.T) {
	localStorage := NewStorage()
	assert.Implements(t, (*MetricsStorage)(nil), localStorage)
}
