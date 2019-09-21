package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	graphql "github.com/graphql-go/graphql"
	godotenv "github.com/joho/godotenv"
)

var schema graphql.Schema

// Start the server
func (server ServerConfig) Start() {
	fmt.Println("\x1b[1mStarting Server:\x1b[0m\n")

	// Load env variables
	fmt.Print("  - Loading \x1b[1m\x1b[96mEnvironment Vairables\x1b[0m from \x1b[3m.env\x1b[0m")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Println("\n    [Environment]")
		fmt.Printf("    - ENV_FILE: %s\n", envFile)
		fmt.Printf("    - DEV: %t\n", os.Getenv("DEV") == "true")
		fmt.Printf("    - VERBOSE: %t\n", isVerbose)
		fmt.Printf("    - STARTUP_ONLY: %t\n", os.Getenv("STARTUP_ONLY") == "true")
		fmt.Printf("    \x1b[92mLoaded\x1b[92m")
	}
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("        Error loading '%s': %v\n\n", envFile, err)
	} else {
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	}
	if isVerbose {
		fmt.Println()
	}

	// Connecting to Database
	fmt.Printf("  - Connecting to \x1b[1m\x1b[94m%s\x1b[0m", server.DatabaseDriver.Name)
	err = ConnectDatabase(server.DatabaseDriver)
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Failed to create connect to the database, error: %v\n\n", err)
		os.Exit(1)
	}
	if isVerbose {
		fmt.Printf("\n    \x1b[92mConnected\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	if isVerbose {
		fmt.Println()
	}

	// Create the GraphQL Schema
	fmt.Print("  - Creating the \x1b[1m\x1b[95mGraphQL Schema\x1b[0m from models")
	schema, err = CreateSchema(server)
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Failed to create GraphQL schema, error: %v\n\n", err)
		os.Exit(1)
	}
	if isVerbose {
		fmt.Printf("    \x1b[92mCreated\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	if isVerbose {
		fmt.Println()
	}

	// Checking Server Config
	fmt.Print("  - Verifying the \x1b[1m\x1b[93mServerConfig\x1b[0m")
	port := server.Port
	if port < 0 || port > 65535 {
		err = fmt.Errorf("Port (%d) must be between 0 and 65535", port)
	}
	endpoint := server.Endpoint
	if endpoint != "" && !strings.HasPrefix(endpoint, "/") {
		err = fmt.Errorf("Endpoint (%s) must start with a '/'", endpoint)
	}
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m\n")
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(1)
	}
	if isVerbose {
		fmt.Printf("\n    \x1b[92mVerified\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")

	// Start the server
	if os.Getenv("STARTUP_ONLY") == "true" {
		fmt.Println()
		return
	}
	startWebServer(server, schema)
}

func startWebServer(server ServerConfig, s graphql.Schema) {
	schema = s
	endpoint := server.Endpoint
	if endpoint == "" {
		endpoint = "/graphql"
	}
	port := server.Port
	if port == 0 {
		port = 8000
	}
	handler := http.HandlerFunc(graphqlHandler)
	http.Handle(endpoint, requestLogger(methodFilter(allowCors(handler))))

	if isDev {
		endpoint := fmt.Sprintf("http://localhost:%d%s", port, endpoint)
		fmt.Printf("  - Starting at \x1b[1m%s\x1b[0m\n", endpoint)
	}

	fmt.Println("\n\x1b[1mLogs:\x1b[0m\n")
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		nil,
	))
}
