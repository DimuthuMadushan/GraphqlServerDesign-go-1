package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var schema graphql.Schema

func getEmployee(args map[string]interface{}) map[string]interface{} {
	id := fmt.Sprint(args["id"])
	data := make(map[string]interface{})
	link := "http://localhost:9090/employeemgt/employee/" + id
	response, err := http.Get(link)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)

	} else {
		content, _ := ioutil.ReadAll(response.Body)
		responseString := string(content)
		temp := fmt.Sprint(responseString)
		err = json.Unmarshal([]byte(temp), &data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			return data
		}
	}

	return nil
}

func buildSchema() error {
	var employeeType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Employee",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"age": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"employee": &graphql.Field{
					Type: employeeType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return getEmployee(p.Args), nil
					},
				},
			},
		})

	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return nil
}

func main() {
	err := buildSchema()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("Now server is running on 'http://localhost:8080/graphql'")
	http.ListenAndServe(":8080", nil)

}

// func executeQuery(query string, schema graphql.Schema) *graphql.Result {
// 	result := graphql.Do(graphql.Params{
// 		Schema:        schema,
// 		RequestString: query,
// 	})
// 	if len(result.Errors) > 0 {
// 		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
// 	}
// 	return result
// }
