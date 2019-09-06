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
	Models              map[string]Model
	Queries             map[string]Resolvable
	DatabaseDriver      DatabaseDriver
}

// Model -
type Model struct {
	Name       string
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
	FromJSON    func(jsonValue) interface{}
}

// Resolvable -
type Resolvable struct {
}

// DatabaseDriver -
type DatabaseDriver struct {
	Name    string
	Connect func()
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (rec *statusWriter) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
