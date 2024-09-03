package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Incident represents an incident record
type Incident struct {
	IncidentID   string
	Title        string
	Description  string
	Status       string
	CreationDate string
}

// Check if the table exists
func tableExists(client *dynamodb.Client) (bool, error) {
	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create the table if it doesn't exist
func createTable(client *dynamodb.Client) error {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("IncidentID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("IncidentID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("Table created successfully.")
	return nil
}

// Create an incident
func createIncident(client *dynamodb.Client, incident Incident) error {
	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"IncidentID":   &types.AttributeValueMemberS{Value: incident.IncidentID},
			"Title":        &types.AttributeValueMemberS{Value: incident.Title},
			"Description":  &types.AttributeValueMemberS{Value: incident.Description},
			"Status":       &types.AttributeValueMemberS{Value: incident.Status},
			"CreationDate": &types.AttributeValueMemberS{Value: incident.CreationDate},
		},
	})
	return err
}

// Get an incident by ID
func getIncident(client *dynamodb.Client, incidentID string) (*Incident, error) {
	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"IncidentID": &types.AttributeValueMemberS{Value: incidentID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("incident not found")
	}

	incident := Incident{
		IncidentID:   incidentID,
		Title:        result.Item["Title"].(*types.AttributeValueMemberS).Value,
		Description:  result.Item["Description"].(*types.AttributeValueMemberS).Value,
		Status:       result.Item["Status"].(*types.AttributeValueMemberS).Value,
		CreationDate: result.Item["CreationDate"].(*types.AttributeValueMemberS).Value,
	}

	return &incident, nil
}

// Update an incident's status
func updateIncident(client *dynamodb.Client, incidentID, newStatus string) error {
	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"IncidentID": &types.AttributeValueMemberS{Value: incidentID},
		},
		UpdateExpression: aws.String("SET #status = :newStatus"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newStatus": &types.AttributeValueMemberS{Value: newStatus},
		},
	})
	return err
}

// Delete an incident by ID
func deleteIncident(client *dynamodb.Client, incidentID string) error {
	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"IncidentID": &types.AttributeValueMemberS{Value: incidentID},
		},
	})
	return err
}
