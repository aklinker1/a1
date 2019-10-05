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

		Types:  customTypes,
		Enums:  customEnums,
		Models: models,
		DataLoaders: a1.DataLoaderMap{
			"PostgreSQL": cache.CreateDriver(cachedData),
			"MongoDB":    cache.CreateDriver(cachedData),
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
		"_id":     1,
		"message": "Todo 1",
		"user_id": 1,
	},
	2: map[string]interface{}{
		"_id":     2,
		"message": "Todo 2",
		"user_id": 1,
	},
	3: map[string]interface{}{
		"_id":     3,
		"message": "Todo 3",
		"user_id": 2,
	},
	4: map[string]interface{}{
		"_id":     4,
		"message": "Todo 4",
		"user_id": -1,
	},
}

var users = map[interface{}]map[string]interface{}{
	1: map[string]interface{}{
		"_id":       1,
		"username":  "User1",
		"email":     "user1@gmail.com",
		"firstName": "Aaron",
		"lastName":  "Klinker",
	},
	2: map[string]interface{}{
		"_id":       2,
		"username":  "User2",
		"email":     "user2@outlook.com",
		"firstName": "Isaiah",
		"lastName":  "Walker",
	},
}
