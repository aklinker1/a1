package pkg

import (
	"fmt"
	"os"
	"strings"

	graphql "github.com/graphql-go/graphql"
	godotenv "github.com/joho/godotenv"
)

var schema graphql.Schema

// Start the server
func Start(serverConfig ServerConfig) {
	fmt.Println("\x1b[1mStarting Server:\x1b[0m\n")

	// Load env variables
	fmt.Print("  - Loading \x1b[1m\x1b[96mEnvironment Vairables\x1b[0m")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("        Error loading '%s': %v\n\n", envFile, err)
	} else {
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	}

	// Create the GraphQL Schema
	fmt.Print("  - Creating the \x1b[1m\x1b[95mGraphQL Schema\x1b[0m")
	schema, err = CreateSchema(serverConfig)
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Failed to create GraphQL schema, error: %v\n\n", err)
		os.Exit(1)
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")

	// Connecting to Database
	fmt.Printf("  - Connecting to \x1b[1m\x1b[94m%s\x1b[0m", serverConfig.DatabaseDriver.Name)
	err = ConnectDatabase(serverConfig.DatabaseDriver)
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Failed to create connect to the database, error: %v\n\n", err)
		os.Exit(1)
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")

	// Checking Server Config
	fmt.Print("  - Verifying the \x1b[1m\x1b[93mServerConfig\x1b[0m")
	port := serverConfig.Port
	if port < 0 || port > 65535 {
		err = fmt.Errorf("Port (%d) must be between 0 and 65535", port)
	}
	endpoint := serverConfig.Endpoint
	if endpoint != "" && !strings.HasPrefix(endpoint, "/") {
		err = fmt.Errorf("Endpoint (%s) must start with a '/'", endpoint)
	}
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(1)
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")

	// Start the server
	StartWebServer(serverConfig, schema)
}
