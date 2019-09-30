package new

type contextKey string

const (
	// ContextKeyAuthHeader is used to store the "authorization" header from the request into the graphql context
	ContextKeyAuthHeader contextKey = "authKey"
)
