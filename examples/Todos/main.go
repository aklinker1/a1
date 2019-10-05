package main

import (
	a1 "github.com/aklinker1/a1/pkg"
)

func main() {
	server := a1.ServerConfig{
		Types:       customTypes,
		Enums:       customEnums,
		Models:      models,
		DataLoaders: dataLoaders,

		EnableIntrospection: true,
	}
	a1.Start(server)
}
