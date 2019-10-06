package pkg

import (
	"fmt"
	"net/http"
	"strings"

	graphql "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// DataMap -
type DataMap = map[string]interface{}

// StringMap -
type StringMap = map[string]string

// RequestedField -
type RequestedField struct {
	Field       FinalGraphQLField
	WasRequired bool
	InnerFields interface{}
	Model       *FinalModel
}

// RequestedFieldMap -
type RequestedFieldMap = map[string]RequestedField

// DataLoader -
type DataLoader struct {
	Connect     func() error
	GetOne      func(model *FinalModel, args DataMap, fields RequestedFieldMap) (DataMap, error)
	GetMultiple func(model *FinalModel, args DataMap, fields RequestedFieldMap) ([]DataMap, error)
	Update      func(model *FinalModel, args DataMap, inputData DataMap, fields RequestedFieldMap) (DataMap, error)
}

// DataLoaderMap -
type DataLoaderMap = map[string]DataLoader

// FinalDataLoader -
type FinalDataLoader struct {
	Name        string
	Connect     func() error
	GetOne      func(model *FinalModel, args DataMap, fields RequestedFieldMap) (DataMap, error)
	GetMultiple func(model *FinalModel, args DataMap, fields RequestedFieldMap) ([]DataMap, error)
}

// FinalDataLoaderMap -
type FinalDataLoaderMap = map[string]*FinalDataLoader

// DataLoaderConfig -
type DataLoaderConfig struct {
	Source string
	Group  string
}

// FinalDataLoaderConfig -
type FinalDataLoaderConfig struct {
	DataLoader *FinalDataLoader
	Source     string
	Group      string
}

// ASTValue - Raw value from graphql framework
type ASTValue = ast.Value

// CustomType -
type CustomType struct {
	Description string
	ToJSON      func(value interface{}) interface{}
	FromJSON    func(value interface{}) interface{}
	FromLiteral func(value ASTValue) interface{}
}

// CustomTypeMap -
type CustomTypeMap = map[string]CustomType

// FinalCustomType -
type FinalCustomType struct {
	Name        string
	GraphQLType graphql.Type
	Description string
	ToJSON      func(value interface{}) interface{}
	FromJSON    func(value interface{}) interface{}
	FromLiteral func(value ASTValue) interface{}
}

// FinalCustomTypeMap -
type FinalCustomTypeMap = map[string]*FinalCustomType

// EnumValue -
type EnumValue struct {
	Value             interface{}
	Description       string
	DeprecationReason string
}

// EnumValueMap -
type EnumValueMap = map[string]EnumValue

// Enum -
type Enum struct {
	Description string
	Values      EnumValueMap
}

// EnumMap -
type EnumMap = map[string]Enum

// LinkType -
type LinkType = int

const (
	// OneToOne -
	OneToOne LinkType = 0
	// OneToMany -
	OneToMany LinkType = 1
	// ManyToOne -
	ManyToOne LinkType = 2
)

// GraphQLField -
type GraphQLField = interface{}

// FieldMap -
type FieldMap = map[string]GraphQLField

// FinalGraphQLField -
type FinalGraphQLField = interface{}

// FinalFieldMap -
type FinalFieldMap = map[string]FinalGraphQLField

// Field -
type Field struct {
	Description       string
	DeprecationReason string
	Type              string
	Hidden            bool
	PrimaryKey        bool
	DataField         string
}

// FinalField -
type FinalField struct {
	Name              string
	Description       string
	DeprecationReason string
	TypeName          string
	Type              *FinalCustomType
	Hidden            bool
	PrimaryKey        bool
	DataField         string
}

// VirtualField -
type VirtualField struct {
	Description       string
	DeprecationReason string
	Type              string
	RequiredFields    []string
	Compute           func(data DataMap) (interface{}, error)
}

// FinalVirtualField -
type FinalVirtualField struct {
	Name              string
	Description       string
	DeprecationReason string
	TypeName          string
	Type              *FinalCustomType
	RequiredFields    []string
	Compute           func(data DataMap) (interface{}, error)
}

// LinkedField -
type LinkedField struct {
	Description       string
	DeprecationReason string
	CustomModelName   string
	LinkedModel       string
	Type              LinkType
	Field             string
	LinkedField       string
}

func (field LinkedField) getCustomModelName(parentModelName string) string {
	name := field.CustomModelName
	if name == "" {
		parentName := parentModelName
		// Change "User.Preferences" to just "Preferences"
		if strings.Contains(parentName, ChildModelDeliminator) {
			parentName = strings.Split(parentName, ChildModelDeliminator)[1]
		}
		name = fmt.Sprintf("%s_%s", parentName, field.LinkedModel)
	}
	return name
}

// FinalLinkedField -
type FinalLinkedField struct {
	Name              string
	Description       string
	DeprecationReason string
	LinkedModelName   string
	LinkedModel       *FinalModel
	Type              LinkType
	Field             string
	LinkedField       string
}

// GraphQLConfig -
type GraphQLConfig struct {
	DisableSelectOne      bool
	DisableSelectMultiple bool
	DisableCreate         bool
	DisableUpdate         bool
	DisableDelete         bool
}

// FinalGraphQLConfig -
type FinalGraphQLConfig struct {
	DisableGetOne      bool
	DisableGetMultiple bool
	DisableCreate      bool
	DisableUpdate      bool
	DisableDelete      bool
}

// Model -
type Model struct {
	Extends     string
	Description string
	DataLoader  DataLoaderConfig
	Fields      FieldMap
	GraphQL     GraphQLConfig
}

// ModelMap -
type ModelMap = map[string]Model

// FinalModel -
type FinalModel struct {
	Name        string
	Extended    string
	Description string
	DataLoader  FinalDataLoaderConfig
	Fields      FinalFieldMap
	DataFields  StringMap
	GraphQL     FinalGraphQLConfig
	PrimaryKey  string
}

func (child FinalModel) extends(parent *FinalModel) *FinalModel {
	// Data loader
	source := parent.DataLoader.Source
	if child.DataLoader.Source != "" {
		source = child.DataLoader.Source
	}
	group := parent.DataLoader.Group
	if child.DataLoader.Group != "" {
		group = child.DataLoader.Group
	}
	dataLoader := parent.DataLoader.DataLoader
	if child.DataLoader.DataLoader != nil {
		dataLoader = child.DataLoader.DataLoader
	}
	dataLoaderConfig := FinalDataLoaderConfig{
		Source:     source,
		Group:      group,
		DataLoader: dataLoader,
	}

	// Fields
	fields := FinalFieldMap{}
	for fieldName, field := range parent.Fields {
		fields[fieldName] = field
	}
	for fieldName, field := range child.Fields {
		fields[fieldName] = field
	}

	// Data Fields
	dataFields := StringMap{}
	for fieldName, field := range fields {
		if regularField, isRegularField := field.(*FinalField); isRegularField {
			dataFields[regularField.DataField] = fieldName
		}
	}

	// Basic details
	name := parent.Name
	if child.Name != "" {
		name = child.Name
	}
	description := parent.Description
	if child.Description != "" {
		description = child.Description
	}
	primaryKey := parent.PrimaryKey
	if child.PrimaryKey != "" {
		primaryKey = child.PrimaryKey
	}
	return &FinalModel{
		Name:        name,
		Extended:    child.Name,
		Description: description,
		DataLoader:  dataLoaderConfig,
		Fields:      fields,
		DataFields:  dataFields,
		GraphQL:     child.GraphQL,
		PrimaryKey:  primaryKey,
	}
}

// FinalModelMap -
type FinalModelMap = map[string]*FinalModel

// ServerConfig -
type ServerConfig struct {
	EnableIntrospection bool
	Port                int
	Endpoint            string
	Models              ModelMap
	DataLoaders         func() DataLoaderMap
	Types               CustomTypeMap
	Enums               EnumMap
}

// FinalServerConfig -
type FinalServerConfig struct {
	EnableIntrospection bool
	Port                int
	Endpoint            string
	Models              FinalModelMap
	DataLoaders         FinalDataLoaderMap
	Types               FinalCustomTypeMap

	GraphQLSchema    graphql.Schema
	GraphQLQueries   graphql.Fields
	GraphQLMutations graphql.Fields
	GraphqlTypes     []graphql.Type
}

type requestBody struct {
	Query     string  `json:"query"`
	Variables DataMap `json:"variables"`
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	type ArgumentMap = DataMap
}

// Argument -
type Argument struct {
	Name         string
	Description  string
	Type         string
	DefaultValue interface{}
}

// Resolvable -
type Resolvable struct {
	Name         string
	Model        *FinalModel
	Description  string
	Arguments    []Argument
	ResturnsList bool
	Resolver     func(args DataMap, fields RequestedFieldMap) (interface{}, error)
}
