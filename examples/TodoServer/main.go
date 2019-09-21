package main

import (
	a1 "github.com/aklinker1/a1/pkg"
	cache "github.com/aklinker1/a1/pkg/drivers/cache"
)

func main() {
	server := a1.ServerConfig{
		EnableIntrospection: true,
		Port:                8000,
		Endpoint:            "/graphql",

		Models: []a1.Model{
			todoModel,
			userModel,
		},
		DatabaseDriver: cache.CreateDriver(cachedData),
	}
	server.Start()
}

// data

var cachedData = map[string]map[interface{}]interface{}{
	"todos": todos,
	"users": users,
}

var todos = map[interface{}]interface{}{
	1: map[interface{}]interface{}{
		"id":      1,
		"title":   "Todo 1",
		"_userId": 1,
	},
	2: map[interface{}]interface{}{
		"id":      2,
		"title":   "Todo 2",
		"_userId": 1,
	},
	3: map[interface{}]interface{}{
		"id":      3,
		"title":   "Todo 3",
		"_userId": 2,
	},
	4: map[interface{}]interface{}{
		"id":      4,
		"title":   "Todo 4",
		"_userId": -1,
	},
}

var users = map[interface{}]interface{}{
	1: map[interface{}]interface{}{
		"id":       1,
		"username": "User1",
	},
	2: map[interface{}]interface{}{
		"id":    2,
		"title": "User2",
	},
}

// Models

var todoModel = a1.Model{
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
		"_userId": a1.Field{
			Type: "User",
			Linking: a1.FieldLink{
				ForeignKey: "id",
				Type:       a1.ManyToOne,
			},
		},
	},
}

var userModel = a1.Model{
	Name:       "User",
	Table:      "users",
	PrimaryKey: "id",
	Fields: map[string]a1.Field{
		"id": a1.Field{
			Type: "Int",
		},
		"username": a1.Field{
			Type: "String",
		},
	},
}
