package pkg

import (
	"strings"
)

func selectOneQuery(model Model, databaseDriver DatabaseDriver) *Resolvable {
	return &Resolvable{
		Model:   model,
		Name:    strings.ToLower(model.Name),
		Returns: model,
		Arguments: []Argument{
			Argument{
				Name: model.PrimaryKey,
				Type: "Int",
			},
		},
		Resolver: func(args StringMap, fields StringMap) (StringMap, error) {
			return databaseDriver.SelectOne(model, args[model.PrimaryKey], fields)
		},
	}
}
