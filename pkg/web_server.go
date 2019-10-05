package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aklinker1/a1/pkg/utils"
	graphql "github.com/graphql-go/graphql"
)

var serverConfig *FinalServerConfig

func startWebServer(finalServerConfig *FinalServerConfig) {
	serverConfig = finalServerConfig
	endpoint := serverConfig.Endpoint
	if endpoint == "" {
		endpoint = "/graphql"
	}
	port := serverConfig.Port
	if port == 0 {
		port = 8000
	}
	handler := http.HandlerFunc(graphqlHandler)
	http.Handle(endpoint, requestLogger(methodFilter(allowCors(handler))))

	if isDev {
		portPlusEndpoint := fmt.Sprintf("%d%s", port, endpoint)
		outboundIP, ipError := utils.GetOutboundIP()
		fmt.Printf("  - Starting the \x1b[1mWeb Server\x1b[0m")
		if utils.IsVerbose() {
			utils.Log("")
			utils.LogWhite("[Server Config]")
			utils.Log("  - Introspection: %t", finalServerConfig.EnableIntrospection)
			utils.Log("  - Endpoint:      %s", endpoint)
			utils.Log("  - Port:          %d", port)
			fmt.Printf("    \x1b[92mStarted\x1b[92m")
		}
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
		fmt.Println()
		fmt.Println("\x1b[1mServer started at:\x1b[0m")
		fmt.Printf("  - Device:  \x1b[4mhttp://localhost:%s\x1b[0m\n", portPlusEndpoint)
		if ipError == nil {
			fmt.Printf("  - Network: \x1b[4mhttp://%s:%s\x1b[0m\n", outboundIP, portPlusEndpoint)
		}
	}

	fmt.Println("\n\x1b[1mLogs:\x1b[0m")
	fmt.Println()
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		nil,
	))
}

func graphqlHandler(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("Failed to read body: %v", err)))
		return
	}

	request := requestBody{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(fmt.Sprintf("Failed to parse body JSON: %v", err)))
		return
	}

	ctx := context.WithValue(context.Background(), ContextKeyHeader, req.Header.Get("authorization"))

	params := graphql.Params{
		Schema:         serverConfig.GraphQLSchema,
		RequestString:  request.Query,
		Context:        ctx,
		VariableValues: request.Variables,
	}

	// Execute the request
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		var finalStatus int
		if result.Data != nil {
			responseMap := result.Data.(DataMap)
			for _, value := range responseMap {
				if status, ok := value.(int); ok {
					if status > finalStatus {
						finalStatus = status
					}
				}
			}
		}
		if finalStatus != 0 {
			res.WriteHeader(finalStatus)
		} else {
			res.WriteHeader(500)
		}
	}
	json.NewEncoder(res).Encode(result)
}
