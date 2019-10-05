package main

import (
	a1 "github.com/aklinker1/a1/pkg"
	cache "github.com/aklinker1/a1/pkg/data_loaders/cache"
)

func main() {
	server := a1.ServerConfig{
		EnableIntrospection: true,

		Types:  customTypes,
		Enums:  customEnums,
		Models: models,
		DataLoaders: a1.DataLoaderMap{
			"PostgreSQL": cache.CreateDataLoader(cachedData),
		},
	}
	a1.Start(server)
}
