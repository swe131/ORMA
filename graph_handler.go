// handlers/graph_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"orma/db"
	"orma/models"

	"go.mongodb.org/mongo-driver/bson"
)

// GetGraphDataHandler handles the endpoint for fetching graph data.
func GetGraphDataHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from authentication token or session
	userID := getUserIDFromRequest(r)

	// Define the pipeline to group and sum the amounts by month for subscriptions
	subscriptionPipeline := []bson.D{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "month", Value: bson.D{{Key: "$month", Value: "$date"}}},
			}},
			{Key: "totalAmount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "month", Value: "$_id.month"},
			{Key: "totalAmount", Value: 1},
		}}},
	}

	// Aggregate data using pipeline for subscriptions
	subscriptionCursor, err := db.SpendCollection.Aggregate(db.Ctx, subscriptionPipeline)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error aggregating subscription data: %v", err), http.StatusInternalServerError)
		return
	}
	defer subscriptionCursor.Close(db.Ctx)

	// Store aggregated data in graph collection for subscriptions
	var subscriptionGraphData []models.Graph
	for subscriptionCursor.Next(db.Ctx) {
		var result models.Graph
		err := subscriptionCursor.Decode(&result)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding subscription cursor data: %v", err), http.StatusInternalServerError)
			return
		}
		subscriptionGraphData = append(subscriptionGraphData, result)
	}

	response, err := json.Marshal(subscriptionGraphData)
	if err != nil {
		http.Error(w, "Failed to serialize graph data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
