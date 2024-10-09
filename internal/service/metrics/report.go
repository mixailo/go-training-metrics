package metrics

type Type string

const (
	TypeCounter Type = "counter"
	TypeGauge   Type = "gauge"
)

func (t Type) String() string {
	return string(t)
}

type Data struct {
	CounterType Type
	Name        string
	Value       string
}

type Report struct {
	value map[string]Data
}

func NewReport() Report {
	r := Report{}
	r.value = make(map[string]Data)
	return r
}

func (r *Report) All() []Data {
	result := make([]Data, 0, len(r.value))
	for _, v := range r.value {
		result = append(result, v)
	}

	return result
}

func (r *Report) Get(name string) (result Data, ok bool) {
	result, ok = r.value[name]
	return
}

func (r *Report) Has(name string) bool {
	_, ok := r.value[name]
	return ok
}

func (r *Report) Add(counterType Type, name, value string) {
	r.value[name] = Data{CounterType: counterType, Name: name, Value: value}
}
