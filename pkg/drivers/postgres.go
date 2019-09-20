package drivers

import (
	"fmt"

	framework "github.com/aklinker1/graphql-framework/pkg"
)

var todos = map[interface{}]string{
	1: "Todo 1",
	2: "Todo 2",
	3: "Todo 3",
}

// CreatePostgresDriver -
func CreatePostgresDriver() framework.DatabaseDriver {
	return framework.DatabaseDriver{
		Name:    "PostgreSQL",
		Connect: func() {},
		SelectOne: func(model framework.Model, id interface{}, fields framework.StringMap) (framework.StringMap, error) {
			fmt.Println("Postgres.selectOne(", model, id, fields, ")")
			output := map[string]interface{}{}
			output["id"] = id
			output["title"] = todos[id]
			fmt.Println("todo", output)
			return output, nil
		},
	}
}
