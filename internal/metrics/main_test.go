package metrics

import (
	"github.com/stretchr/testify/require"
	"testing"
)
import (
	"github.com/stretchr/testify/assert"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			t:    TypeGauge,
			want: "gauge",
		},
		{
			t:    TypeCounter,
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.t.String())
		})
	}
}

func TestNewReport(t *testing.T) {
	report := NewReport()
	require.NotNil(t, report)
	assert.IsType(t, report, Report{})
}

func TestReport_Add(t *testing.T) {
	report := NewReport()
	assert.False(t, report.Has("test"))
	report.Add(TypeCounter, "test", "value")
	assert.True(t, report.Has("test"))
}

func TestReport_All(t *testing.T) {
	report := NewReport()
	report.Add(TypeCounter, "test1", "value")
	report.Add(TypeCounter, "test2", "value")
	report.Add(TypeCounter, "test3", "value")
	report.Add(TypeCounter, "test3", "value")

	assert.Equal(t, len(report.All()), 3)
	assert.True(t, report.Has("test1"))
	assert.True(t, report.Has("test2"))
	assert.True(t, report.Has("test3"))
	assert.False(t, report.Has("test4"))
}

func TestReport_Get(t *testing.T) {
	report := NewReport()
	report.Add(TypeCounter, "test1", "value")

	tests := []struct {
		name  string
		field string
		value string
		found bool
		equal bool
	}{
		{
			name:  "Existing",
			field: "test1",
			value: "value",
			found: true,
			equal: true,
		},
		{
			name:  "Invalid value",
			field: "test1",
			value: "invalid",
			found: true,
			equal: false,
		},
		{
			name:  "Inexistent",
			field: "test2",
			value: "invalid",
			found: false,
			equal: false,
		},
	}

	for _, tt := range tests {
		value, ok := report.Get(tt.field)
		if tt.found == true {
			assert.True(t, ok)
			if tt.equal {
				assert.Equal(t, tt.value, value.Value)
			}
		} else {
			assert.False(t, ok)
		}
	}
}

func TestReport_Has(t *testing.T) {
	report := NewReport()
	report.Add(TypeCounter, "test1", "value")

	tests := []struct {
		name  string
		field string
		value string
		found bool
	}{
		{
			name:  "Existing",
			field: "test1",
			value: "value",
			found: true,
		},
		{
			name:  "Inexistent",
			field: "test2",
			value: "invalid",
			found: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.found, report.Has(tt.field))
	}
}
