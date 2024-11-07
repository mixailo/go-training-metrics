package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	t.Run("new report creation", func(t *testing.T) {
		require.NotNil(t, report)
		assert.IsType(t, Report{}, report)
	})
}

func TestReport_Add(t *testing.T) {
	report := NewReport()
	assert.False(t, report.Has("test"))
	val := 10.3
	metrics := Metrics{ID: "test", MType: TypeGauge.String(), Value: &val}
	report.Add(metrics)
	t.Run("add gauge", func(t *testing.T) {
		assert.True(t, report.Has("test"), "existing")
		assert.False(t, report.Has("test2"), "non-existing")
	})
}

func TestReport_Length(t *testing.T) {
	report := NewReport()
	report.AddUnConverted(TypeCounter, "test1", "1")
	report.AddUnConverted(TypeCounter, "test2", "2")

	t.Run("length", func(t *testing.T) {
		assert.Equal(t, 2, report.Length())
	})
}

func TestReport_All(t *testing.T) {
	report := NewReport()
	report.AddUnConverted(TypeCounter, "test1", "1")
	report.AddUnConverted(TypeCounter, "test2", "2")
	report.AddUnConverted(TypeCounter, "test3", "3")
	report.AddUnConverted(TypeCounter, "test3", "4")

	t.Run("all", func(t *testing.T) {
		assert.Equal(t, 3, report.Length())
		assert.True(t, report.Has("test1"))
		assert.True(t, report.Has("test2"))
		assert.True(t, report.Has("test3"))
		assert.False(t, report.Has("test4"))
	})
}

func TestReport_Get(t *testing.T) {
	report := NewReport()
	report.AddUnConverted(TypeCounter, "test1", "10")

	tests := []struct {
		name  string
		field string
		value int
		found bool
		equal bool
	}{
		{
			name:  "Existing",
			field: "test1",
			value: 10,
			found: true,
			equal: true,
		},
		{
			name:  "Invalid value",
			field: "test1",
			value: 11,
			found: true,
			equal: false,
		},
		{
			name:  "Inexistent",
			field: "test2",
			value: 11,
			found: false,
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := report.Get(tt.field)
			if tt.found == true {
				assert.True(t, ok)
				if tt.equal {
					if value.Delta == nil {
						t.Fatal("nil value (probably value not parsed)")
					} else {
						assert.Equal(t, int64(tt.value), *value.Delta)
					}
				}
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestReport_Has(t *testing.T) {
	report := NewReport()
	report.AddUnConverted(TypeCounter, "test1", "10")

	tests := []struct {
		name  string
		field string
		found bool
	}{
		{
			name:  "Existing",
			field: "test1",
			found: true,
		},
		{
			name:  "Inexistent",
			field: "test2",
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.found, report.Has(tt.field))
		})
	}
}

func TestNewMetrics(t *testing.T) {
	var gaugeVal float64
	var counterVal int64
	gaugeVal = 10
	counterVal = 20

	type args struct {
		counterType Type
		name        string
		value       string
	}
	tests := []struct {
		name    string
		args    args
		wantRes Metrics
		wantErr bool
	}{
		{
			name: "successful creation of a gauge",
			args: args{
				counterType: TypeGauge,
				name:        "gauge",
				value:       "10",
			},
			wantRes: Metrics{
				ID:    "gauge",
				MType: TypeGauge.String(),
				Value: &gaugeVal,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful creation of a gauge",
			args: args{
				counterType: TypeGauge,
				name:        "gauge",
				value:       "jj",
			},
			wantRes: Metrics{
				ID:    "gauge",
				MType: TypeGauge.String(),
				Value: &gaugeVal,
			},
			wantErr: true,
		},
		{
			name: "successful creation of a counter",
			args: args{
				counterType: TypeCounter,
				name:        "counter",
				value:       "20",
			},
			wantRes: Metrics{
				ID:    "counter",
				MType: TypeCounter.String(),
				Delta: &counterVal,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful creation of a counter",
			args: args{
				counterType: TypeCounter,
				name:        "counter",
				value:       "ggg",
			},
			wantRes: Metrics{
				ID:    "counter",
				MType: TypeCounter.String(),
				Delta: &counterVal,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := NewMetrics(tt.args.counterType, tt.args.name, tt.args.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.EqualValues(t, tt.wantRes, gotRes)
			}
		})
	}
}
