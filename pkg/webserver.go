package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	graphql "github.com/graphql-go/graphql"
)

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
			responseMap := result.Data.(StringMap)
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
