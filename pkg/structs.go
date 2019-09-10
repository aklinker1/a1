package pkg

import "net/http"

type requestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// ServerConfig -
type ServerConfig struct {
	EnableIntrospection bool
	Port                int
	Endpoint            string
	Models              []Model
	Queries             []Resolvable
	DatabaseDriver      DatabaseDriver
}

// Model -
type Model struct {
	Name       string
	Table      string
	PrimaryKey string
	Fields     map[string]Field
}

// Field -
type Field struct {
	Type Scalar
}

// Scalar -
type Scalar struct {
	Name        string
	Description string
	FromJSON    func(jsonValue interface{}) interface{}
	ToJSON      func(value interface{}) interface{}
}

// Resolvable -
type Resolvable struct {
}

// DatabaseDriver -
type DatabaseDriver struct {
	Name    string
	Connect func()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ModelMapItem -
type ModelMapItem struct {
	Model     Model
	Queries   []Resolvable
	Mutations []Resolvable
}

// ModelMap -
type ModelMap = map[string]ModelMapItem

// CustomScalarMap -
type CustomScalarMap = map[string]Scalar

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
