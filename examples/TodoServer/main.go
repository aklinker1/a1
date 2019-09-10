package main

import (
	framework "github.com/aklinker1/graphql-framework/pkg"
)

func main() {
	server := framework.ServerConfig{
		EnableIntrospection: true,
		Port:                8080,
		Endpoint:            "/graphql",

		Models: []framework.Model{
			framework.Model{
				Name:       "Todo",
				Table:      "todos",
				PrimaryKey: "id",
				Fields: map[string]framework.Field{
					"id": framework.Field{
						Type: framework.Int,
					},
					"title": framework.Field{
						Type: framework.String,
					},
				},
			},
		},
		DatabaseDriver: framework.DatabaseDriver{
			Name:    "Postgres",
			Connect: func() {},
		},
	}
	server.Start()
}
