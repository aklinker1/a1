package postgres

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
)

var todos = map[interface{}]string{
	1: "Todo 1",
	2: "Todo 2",
	3: "",
}

// CreateDriver -
func CreateDriver() a1.DatabaseDriver {
	return a1.DatabaseDriver{
		Name:    "PostgreSQL",
		Connect: func() {},
		SelectOne: func(model a1.Model, id interface{}, fields a1.StringMap) (a1.StringMap, error) {
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
