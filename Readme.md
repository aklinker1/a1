<img width="200" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png">

A1 is a GraphQL framework written in Go. It is designed to remove boilerplate code by defining only the models returned by the API and their relationships. This makes setting up an entire API is easy.

## Getting Started

You can either start from scratch by installing the library directly using `go get`, or you can clone a boilerplat prohject to get up and running faster.

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
    postgres "github.com/aklinker1/graphql-framework/pkg/drivers/postgres"
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
        DatabaseDriver: postgres.CreateDriver(),
    }
    server.Start()
}
```

## FAQ

### 1. Is this an ORM?

No, A1 does not do any database interaction. All database interactions go though the `DatabaseDriver`, while A1 simply tells the driver what it wants done. This also means that A1 does not handle database setup or teardown. You qwill have to create the tables and manage migrations.

### 2. Can I still customize a `selectOne` query or any other queries where I don't want the default behavior?

Of course! Check out [this page]() to find out how to override any default behaviors/

### 3. Do you support subscriptions?

As of now, no.

## Documentation

To checkout the full documentation and examples for getting started, checkout the [`docs/`]() folder.
