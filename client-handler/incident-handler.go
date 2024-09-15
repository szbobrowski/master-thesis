package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"      // Added log package
	"net/http" // Added for stdout
	"time"

	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const grpcAddress = "localhost:50052"

var incidentClient IncidentServiceClient // gRPC client for GraphQL

// Create Incident GraphQL Type
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

// Define the GraphQL query root
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
						return nil, fmt.Errorf("incidentID is required")
					}

					// Call GetIncident method from the gRPC server
					req := &GetIncidentRequest{
						IncidentID: incidentID,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					resp, err := incidentClient.GetIncident(ctx, req)
					if err != nil {
						st, _ := status.FromError(err)
						log.Printf("Failed to retrieve incident with ID: %s, error: %v\n", incidentID, st.Message()) // Log error
						return nil, fmt.Errorf("gRPC error: %v", st.Message())
					}

					log.Printf("Incident retrieved: %+v\n", resp.Incident) // Log success
					return resp.Incident, nil
				},
			},
		},
	},
)

// Define mutation for creating, updating, and deleting an incident
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

					// Call CreateIncident method from the gRPC server
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
						log.Printf("Failed to create incident, error: %v\n", err) // Log error
						return nil, fmt.Errorf("gRPC error: %v", err.Error())
					}

					log.Printf("Incident created: %+v\n", resp.Incident) // Log success
					return resp.Incident, nil
				},
			},

			// Mutation for updating an incident
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

					// Call UpdateIncident method from the gRPC server
					req := &UpdateIncidentRequest{
						IncidentID: incidentID,
						Status:     status,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					resp, err := incidentClient.UpdateIncident(ctx, req)
					if err != nil {
						log.Printf("Failed to update incident with ID: %s, error: %v\n", incidentID, err) // Log error
						return nil, fmt.Errorf("gRPC error: %v", err.Error())
					}

					log.Printf("Incident updated: %+v\n", resp.Incident) // Log success
					return resp.Incident, nil
				},
			},

			// Mutation for deleting an incident
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

					// Call DeleteIncident method from the gRPC server
					req := &DeleteIncidentRequest{
						IncidentID: incidentID,
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()

					_, err := incidentClient.DeleteIncident(ctx, req)
					if err != nil {
						log.Printf("Failed to delete incident with ID: %s, error: %v\n", incidentID, err) // Log error
						return nil, fmt.Errorf("gRPC error: %v", err.Error())
					}

					log.Printf("Incident deleted with ID: %s\n", incidentID) // Log success
					return map[string]interface{}{
						"success": true,
					}, nil
				},
			},
		},
	},
)

// Initialize gRPC client
func initIncidentGrpcClient() {
	// Dial to gRPC server
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	incidentClient = NewIncidentServiceClient(conn)
	log.Println("Connected to gRPC server at", grpcAddress, "for GraphQL") // Log connection to gRPC server
}

// Define the GraphQL schema
var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	},
)

// GraphQL handler
func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Println("Error parsing GraphQL request body") // Log error
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Log the incoming query
	log.Printf("Received GraphQL query: %s\n", params.Query)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: params.Query,
	})

	json.NewEncoder(w).Encode(result)
}
