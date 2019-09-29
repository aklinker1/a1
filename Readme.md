# A1

A1 is a GraphQL framework written in Go. It is designed to remove tedious boilerplate and repetitive code so that the only thing it needs to spin up a server is a list of models and their relationships.

## Getting Started

You can either start from scratch by installing the library directly using `go get`, or you can clone a boilerplate project to get up and running faster.

```bash
# Install the library directly
go get github.com/aklinker1/a1

# Clone a boilerplate project
git clone https://github.com/aklinker1/a1-boilerplate
```

Then simply import the library and your driver of choice, define your models and server config, and call `Start`!

```golang
import (
    framework "github.com/aklinker1/graphql-framework/pkg"
    cache "github.com/aklinker1/graphql-framework/pkg/drivers/cache"
)

func main() {
    server := framework.ServerConfig{
        EnableIntrospection: true,
        Port:                8000,
        Endpoint:            "/graphql",

        Models: []framework.Model{
            framework.Model{
                Name:       "Todo",
                Table:      "todos",
                PrimaryKey: "id",
                Fields: map[string]framework.Field{
                    "id": framework.Field{
                        Type: "Int",
                    },
                    "title": framework.Field{
                        Type: "String",
                    },
                },
            },
        },
        DatabaseDriver: cache.CreateDriver(cachedData),
    }
    server.Start()
}

// This type might be confusing, but it is just a map of tables by their name,
// and each table is a map of items by their primary key.
cachedData := map[string]map[interface{}]map[string]interface{}{
    "todos": map[interface{}]map[string]interface{}{
        1: map[string]interface{}{
            id:    1,
            title: "Todo #1",
        },
        2: interface{}{
            id:    2,
            title: "Some other todo",
        },
    },
}
```

For simplicity, A1 includes a "cache" database driver that can be used for testing out the library, but should not be used in production for caching. All it does is look at a map of tables, and modify/access each table accordingly. Once the server restarts, the tables will be reset.

## Documentation

To checkout the full documentation and examples for getting started, checkout the [`docs/`](https://github.com/aklinker1/a1/tree/master/docs) folder.

## FAQ

### 1. Is this an ORM?

No, A1 does not do any database interaction. All database interactions go though the `DatabaseDriver`, while A1 simply tells the driver what it wants done. This also means that A1 does not handle database setup or teardown. You will have to create the tables and manage migrations yourself.

This also means that any emergent behaviors and relationships, such as `Many to Many`, are not directly supported. To find out how to implement these relationships, checkout the [`linking documentation`](https://github.com/aklinker1/a1/tree/master/docs/linking.md)

### 2. Can I still customize a `selectOne` query or any other queries where I don't want the default behavior?

Of course! Check out [this page](https://github.com/aklinker1/a1/tree/master/docs/extending_behavior.md) to find out how to override any default behaviors, and add custom resolvers independent of the models

### 3. Do you support subscriptions?

No, and as of now there are no plans on doing so.
