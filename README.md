<img width="200" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png"> &emsp;__README__

A1 is a GraphQL framework written in Go. It is designed to remove boilerplate code by defining only the models returned by the API and their relationships. This makes setting up an entire API a piece of cake.

## Installation and Usage

You can either start from scratch by installing the library directly using `go get`, or you can clone a boilerplate project to get up and running faster.

```bash
# Install the library directly
go get github.com/aklinker1/a1

# Clone a boilerplate project
git clone https://github.com/aklinker1/a1-boilerplate.git
```

Then to use the framework, simply create a new `a1.ServerConfig` and then call `a1.Start(ServerConfig)` to start the server.

```go
import (
    a1 "github.com/aklinker1/a1/pkg"
)

func main() {
    server := a1.ServerConfig{
        DataLoaders: a1.DataLoaderMap{ /* define data loaders */ },
        Models:      a1.ModelMap{ /* define models loaders */ },
        Port:        8000,
        Endpoint:    "/graphql",
    }
    a1.Start(server)
}
```

## Basic Example

A simple todo application with users and todos. This example includes all three types of fields: `Field`, `VirtualField`, and `LinkedField`.

```go
package example

import (
    a1 "github.com/aklinker1/a1/pkg"
    a1Types "github.com/aklinker1/a1/pkg/types"
    postgres "github.com/aklinker1/a1-postgresql/pkg"
)

var dataLoaders  = a1.DataLoaderMap{
    "PostgreSQL": postgres.CreateDataLoader("postgresql://postgres:password@localhost:5432/todos_db")
}

var models = a1.ModelMap{
    "User": a1.Model{
        DataLoader: a1.DataLoaderConfig{
            Source: "PostgreSQL",
            Group:  "users",
        }
        Fields: a1.FieldMap{
            "id": a1.Field{
                PrimaryKey: true,
                Type:       "ID",
            },
            "firstName": "String",
            "lastName":  "String",
            "fullName": a1.VirtualField{
                Type:           "String",
                RequiredFields: []string{"firstName", "lastName"},
                Computed: func (data: a1.DataMap) interface{} {
                    return fmt.Sprintf("%s %s", data["firstName"], data["lastName"])
                },
            },
            "todos": a1.LinkedField{
                Model:       "Todo",
                Type:        a1.OneToMany,
                Field:       "id",
                LinkedField: "_userId",
            },
        },
    },
    "Todo": a1.Model{
        DataLoader: a1.DataLoaderConfig{
            Source: "PostgreSQL",
            Group:  "todos",
        },
        Fields: a1.FieldMap{
            "id": a1.Field{
                PrimaryKey: true,
                Type:       "ID",
            },
            "value":      "String",
            "isChecked":  "Boolean",
            "userId":     "ID",
            "user": a1.LinkedField{
                Model:       "User",
                Type:        a1.ManyToOne,
                Field:       "userId",
                LinkedField: "id",
            },
        },
    }
}

func main() {
    server := a1.ServerConfig{
        DataLoaders: dataLoaders,
        Models:      models,
        Port:        8000,
        Endpoint:    "/graphql",
    }
}
```

And there you go! A fully functioning basic todo application GraphQL backend.

> For a bit larger and more complex version, checkout the example [here](https://github.com/aklinker1/a1/tree/master/examples/TodoServer).
