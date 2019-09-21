package pkg

import (
	"net/http"

	graphql "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type requestBody struct {
	Query     string    `json:"query"`
	Variables StringMap `json:"variables"`
}

// ServerConfig -
type ServerConfig struct {
	EnableIntrospection bool
	Port                int
	Endpoint            string
	Models              []Model
	Scalars             []Scalar
	Queries             []Resolvable
	DatabaseDriver      DatabaseDriver
}

// Model -
type Model struct {
	Name        string
	Description string
	Table       string
	PrimaryKey  string
	Fields      map[string]Field
	GraphQL     GraphQLCustomization
}

// GraphQLCustomization -
type GraphQLCustomization struct {
	disableSelectOne bool
	CustomQueries    []*Resolvable
	CustomMutations  []*Resolvable
}

// Field -
type Field struct {
	Name        string
	Description string
	Type        string
}

// Scalar -
type Scalar struct {
	Name        string
	Description string
	serialize   func(value interface{}) interface{}
	parse       func(value interface{}) interface{}
	parseAST    func(valueAST ASTValue) interface{}
}

// Resolvable -
type Resolvable struct {
	Model       Model
	Name        string
	Description string
	Returns     Model
	Arguments   []Argument
	Resolver    func(args StringMap, fields StringMap) (StringMap, error)
}

// DatabaseDriver -
type DatabaseDriver struct {
	Name      string
	Connect   func()
	SelectOne func(model Model, primaryKey interface{}, fieldMap StringMap) (StringMap, error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ModelMapItem -
type ModelMapItem struct {
	Model     Model
	Queries   []*Resolvable
	Mutations []*Resolvable
}

// ModelMap -
type ModelMap = map[string]ModelMapItem

// CustomScalarMap -
type CustomScalarMap = map[string]graphql.Type

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	type ArgumentMap = StringMap
}

// Argument -
type Argument struct {
	Name         string
	Description  string
	Type         string
	DefaultValue interface{}
}

// StringMap -
type StringMap = map[string]interface{}

// ASTValue -
type ASTValue = ast.Value
