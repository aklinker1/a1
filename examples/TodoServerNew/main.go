package main

import (
	cache "github.com/aklinker1/a1/pkg/drivers/cache"
	a1 "github.com/aklinker1/a1/pkg/new"
)

func main() {
	server := a1.ServerConfig{
		EnableIntrospection: true,
		Port:                8000,
		Endpoint:            "/graphql",

		Types:  customTypes,
		Enums:  customEnums,
		Models: models,
		DataLoaders: a1.DataLoaderMap{
			"PostgreSQL": cache.CreateDriver(cachedData),
			"MongDB":     cache.CreateDriver(cachedData),
			"Realm":      cache.CreateDriver(cachedData),
		},
	}
	a1.Start(server)
}

// data

var cachedData = map[string]map[interface{}]map[string]interface{}{
	"todos": todos,
	"users": users,
}

var todos = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"id":      1,
		"title":   "Todo 1",
		"_userId": 1,
	},
	2: map[string]interface{}{
		"id":      2,
		"title":   "Todo 2",
		"_userId": 1,
	},
	3: map[string]interface{}{
		"id":      3,
		"title":   "Todo 3",
		"_userId": 2,
	},
	4: map[string]interface{}{
		"id":      4,
		"title":   "Todo 4",
		"_userId": -1,
	},
}

var users = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"id":       1,
		"username": "User1",
	},
	2: map[string]interface{}{
		"id":       2,
		"username": "User2",
	},
}
