package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const grpcAddress = "localhost:50052"

var incidentClient IncidentServiceClient

var incidentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Incident",
		Fields: graphql.Fields{
			"incidentID": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"creationDate": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"incident": &graphql.Field{
				Type: incidentType,
				Args: graphql.FieldConfigArgument{
					"incidentID": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					incidentID, ok := p.Args["incidentID"].(string)
					if !ok {
						return nil, fmt.Errorf("Wymagane jest pole incidentID")
					}

					req := &GetIncidentRequest{
						IncidentID: incidentID,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					resp, err := incidentClient.GetIncident(ctx, req)
					if err != nil {
						st, _ := status.FromError(err)
						log.Printf("Nie udało się pobrać incydentu o id: %s, błąd: %v\n", incidentID, st.Message())
						return nil, fmt.Errorf("Błąd z serwera incident-notifier: %v", st.Message())
					}

					log.Printf("Pobrano incydent: %+v\n", resp.Incident)
					return resp.Incident, nil
				},
			},
		},
	},
)

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createIncident": &graphql.Field{
				Type: incidentType,
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"creationDate": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					title := p.Args["title"].(string)
					description := p.Args["description"].(string)
					status := p.Args["status"].(string)
					creationDate := p.Args["creationDate"].(string)

					req := &CreateIncidentRequest{
						Title:        title,
						Description:  description,
						Status:       status,
						CreationDate: creationDate,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					resp, err := incidentClient.CreateIncident(ctx, req)
					if err != nil {
						log.Printf("Nie udało się utworzyć incydentu, error: %v\n", err)
						return nil, fmt.Errorf("Błąd z serwera incident-notifier: %v", err.Error())
					}

					log.Printf("Utworzono incydent: %+v\n", resp.Incident)
					return resp.Incident, nil
				},
			},

			"updateIncident": &graphql.Field{
				Type: incidentType,
				Args: graphql.FieldConfigArgument{
					"incidentID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					incidentID := p.Args["incidentID"].(string)
					status := p.Args["status"].(string)

					req := &UpdateIncidentRequest{
						IncidentID: incidentID,
						Status:     status,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					resp, err := incidentClient.UpdateIncident(ctx, req)
					if err != nil {
						log.Printf("Nie udało się zaktualizować incydentu o id: %s, error: %v\n", incidentID, err)
						return nil, fmt.Errorf("Błąd z serwera incident-notifier: %v", err.Error())
					}

					log.Printf("Zaktualizowano incydent: %+v\n", resp.Incident)
					return resp.Incident, nil
				},
			},

			"deleteIncident": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "DeleteResponse",
					Fields: graphql.Fields{
						"success": &graphql.Field{
							Type: graphql.Boolean,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"incidentID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					incidentID := p.Args["incidentID"].(string)

					req := &DeleteIncidentRequest{
						IncidentID: incidentID,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					_, err := incidentClient.DeleteIncident(ctx, req)
					if err != nil {
						log.Printf("Nie udało się usunąć incydentu o id: %s, error: %v\n", incidentID, err)
						return nil, fmt.Errorf("Błąd z serwera incident-notifier: %v", err.Error())
					}

					log.Printf("Usunięto incydent o id: %s\n", incidentID)
					return map[string]interface{}{
						"success": true,
					}, nil
				},
			},
		},
	},
)

func initIncidentGrpcClient() {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Nie udało się połączyć z serwerem incident-notifier: %v", err)
	}
	incidentClient = NewIncidentServiceClient(conn)
}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	},
)

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Println("Błąd podczas dekodowania treści zapytania")
		http.Error(w, "Błąd podczas dekodowania treści zapytania", http.StatusBadRequest)
		return
	}

	log.Printf("Otrzymano kwerendę GraphQL: %s\n", params.Query)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: params.Query,
	})

	json.NewEncoder(w).Encode(result)
}
