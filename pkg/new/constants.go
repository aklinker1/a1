package new

type contextKey string

// ContextKeyHeader is used to store the "authorization" header from the request into the graphql context
const ContextKeyHeader contextKey = "authKey"

// ChildModelDeliminator -
const ChildModelDeliminator string = "_"
