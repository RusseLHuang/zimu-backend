package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/RusseLHuang/zimu-backend/constant"
	"github.com/RusseLHuang/zimu-backend/podcast"
	"github.com/RusseLHuang/zimu-backend/utils"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	graphqlSchema, err := initGraphQLSchema()
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	utils.InitRedis()

	h := handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	configPort := viper.Get("port")
	port := fmt.Sprintf(":%s", configPort)

	http.Handle("/graphql", httpHeaderMiddleware(h))
	http.ListenAndServe(port, nil)
}

func initGraphQLSchema() (graphql.Schema, error) {
	jwtToken := viper.Get("jwt")

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"podcast": &graphql.Field{
					Type: graphql.NewList(podcast.PodcastType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						httpHeader := p.Context.Value(constant.HEADER).(http.Header)
						token := httpHeader.Get(constant.AUTHORIZATION)

						if token != jwtToken {
							return nil, errors.New("Authorization token must be present")
						}

						return podcast.GetAll(), nil
					},
				},
				"collectionSearch": &graphql.Field{
					Type: graphql.NewList(podcast.CollectionType),
					Args: graphql.FieldConfigArgument{
						"keywords": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						httpHeader := p.Context.Value(constant.HEADER).(http.Header)
						token := httpHeader.Get(constant.AUTHORIZATION)

						if token != jwtToken {
							return nil, errors.New("Authorization token must be present")
						}

						keywords := p.Args["keywords"].(string)
						return podcast.Search(keywords), nil
					},
				},
				"collection": &graphql.Field{
					Type: podcast.CollectionType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						httpHeader := p.Context.Value(constant.HEADER).(http.Header)
						token := httpHeader.Get(constant.AUTHORIZATION)

						if token != jwtToken {
							return nil, errors.New("Authorization token must be present")
						}

						id := p.Args["id"].(string)
						return podcast.GetCollection(id), nil
					},
				},
			},
		},
	)

	schemaConfig := graphql.SchemaConfig{Query: queryType}
	schema, err := graphql.NewSchema(schemaConfig)

	return schema, err
}

func httpHeaderMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "header", r.Header)

		next.ContextHandler(ctx, w, r)
	})
}
