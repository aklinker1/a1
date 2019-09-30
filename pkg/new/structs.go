package new

import (
	"fmt"
	"strings"

	graphql "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// DataMap -
type DataMap = map[string]interface{}

// RequestedField -
type RequestedField struct {
	model       *Model
	InnerFields RequestedFieldMap
}

// RequestedFieldMap -
type RequestedFieldMap = map[string]interface{}

// DataLoader -
type DataLoader struct {
	Connect     func() error
	GetOne      func(model FinalModel, primaryKey interface{}, fields RequestedFieldMap) (DataMap, error)
	GetMultiple func(model FinalModel, searchFields DataMap, fields RequestedFieldMap) ([]DataMap, error)
}

// DataLoaderMap -
type DataLoaderMap = map[string]DataLoader

// FinalDataLoader -
type FinalDataLoader struct {
	Name        string
	Connect     func() error
	GetOne      func(model FinalModel, primaryKey interface{}, fields RequestedFieldMap) (DataMap, error)
	GetMultiple func(model FinalModel, searchFields DataMap, fields RequestedFieldMap) ([]DataMap, error)
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
	Requires          []string
	Compute           func(data DataMap) interface{}
}

// FinalVirtualField -
type FinalVirtualField struct {
	Name              string
	Description       string
	DeprecationReason string
	TypeName          string
	Type              *FinalCustomType
	Requires          []string
	Compute           func(data DataMap) interface{}
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
		if strings.ContainsRune(parentName, '.') {
			parentName = strings.Split(parentName, ".")[1]
		}
		name = fmt.Sprintf("%s.%s", parentName, field.LinkedModel)
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
	DisableSelectOne      bool
	DisableSelectMultiple bool
	DisableCreate         bool
	DisableUpdate         bool
	DisableDelete         bool
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
	DataLoaders         DataLoaderMap
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
}
