package main

import (
	framework "github.com/aklinker1/graphql-framework/pkg"
)

func main() {
	framework.Start(framework.ServerConfig{
		EnableIntrospection: true,
		Port:                8080,
		Endpoint:            "/graphql",

		Models: map[string]framework.Model{
			"todos": framework.Model{
				Name:       "Todo",
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
	})
}
