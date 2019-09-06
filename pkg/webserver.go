package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graphql-go/graphql"
)

// StartWebServer - starts the server
func StartWebServer(serverConfig ServerConfig, s graphql.Schema) {
	schema = s
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
		done := make(chan bool)
		go (func() {
			log.Fatal(http.ListenAndServe(
				fmt.Sprintf(":%d", port),
				nil,
			))
		})()
		fmt.Printf("\nServer started at http://localhost:%d%s\n", port, endpoint)
		fmt.Println("\n  ─────\n")
		<-done
	} else {
		log.Fatal(http.ListenAndServe(
			fmt.Sprintf(":%d", port),
			nil,
		))
	}
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

	ctx := context.WithValue(context.Background(), ContextKeyAuthHeader, req.Header.Get("authorization"))

	params := graphql.Params{
		Schema:         schema,
		RequestString:  request.Query,
		Context:        ctx,
		VariableValues: request.Variables,
	}

	// Execute the request
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		var finalStatus int
		if result.Data != nil {
			responseMap := result.Data.(map[string]interface{})
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
