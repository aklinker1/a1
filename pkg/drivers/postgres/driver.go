package drivers

import (
	"fmt"

	framework "github.com/aklinker1/graphql-framework/pkg"
)

var todos = map[interface{}]string{
	1: "Todo 1",
	2: "Todo 2",
	3: "",
}

// CreatePostgresDriver -
func CreatePostgresDriver() framework.DatabaseDriver {
	return framework.DatabaseDriver{
		Name:    "PostgreSQL",
		Connect: func() {},
		SelectOne: func(model framework.Model, id interface{}, fields framework.StringMap) (framework.StringMap, error) {
			output := map[string]interface{}{}
			todo, ok := todos[id]
			if !ok {
				return nil, fmt.Errorf("'%s' with %s=%v does not exist", model.Name, model.PrimaryKey, id)
			}
			output["id"] = id
			output["title"] = todo
			return output, nil
		},
	}
}
