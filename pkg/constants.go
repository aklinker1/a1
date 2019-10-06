package pkg

type contextKey string

// ContextKeyHeader is used to store the "authorization" header from the request into the graphql context
const ContextKeyHeader contextKey = "authKey"

// ChildModelDeliminator -
const ChildModelDeliminator string = "_"

// Custom Types
const (
	// Bool -
	Bool string = "Bool"

	// Int -
	Int string = "Int"

	// Int64 -
	Int64 string = "Int64"

	// Float -
	Float string = "Float"

	// Float64 -
	Float64 string = "Float64"

	// String -
	String string = "String"

	// Date -
	Date string = "Date"
)
