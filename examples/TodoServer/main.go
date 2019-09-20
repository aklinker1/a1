package main

import (
	framework "github.com/aklinker1/graphql-framework/pkg"
	drivers "github.com/aklinker1/graphql-framework/pkg/drivers"
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
		DatabaseDriver: drivers.CreatePostgresDriver(),
	}
	server.Start()
}
