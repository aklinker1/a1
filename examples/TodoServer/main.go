package main

import (
	a1 "github.com/aklinker1/a1/pkg"
	postgres "github.com/aklinker1/a1/pkg/drivers/postgres"
)

func main() {
	server := a1.ServerConfig{
		EnableIntrospection: true,
		Port:                8000,
		Endpoint:            "/graphql",

		Models: []a1.Model{
			a1.Model{
				Name:       "Todo",
				Table:      "todos",
				PrimaryKey: "id",
				Fields: map[string]a1.Field{
					"id": a1.Field{
						Type: "Int",
					},
					"title": a1.Field{
						Type: "String",
					},
				},
			},
		},
		DatabaseDriver: postgres.CreateDriver(),
	}
	server.Start()
}
