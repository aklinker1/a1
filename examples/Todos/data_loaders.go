package main

import (
	"os"

	a1 "github.com/aklinker1/a1/pkg"
	cache "github.com/aklinker1/a1/pkg/data_loaders/cache"
)

func dataLoaders() a1.DataLoaderMap {
	cachedPassword := os.Getenv("CACHE_PASSWORD")
	return a1.DataLoaderMap{
		"PostgreSQL": cache.CreateDataLoader(cachedData, cachedPassword),
	}
}
