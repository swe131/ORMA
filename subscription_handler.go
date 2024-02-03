// handlers/subscription_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"orma/db"
	"orma/models"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// StatusCodeOK represents HTTP status code 200
	StatusCodeOK = http.StatusOK
	// StatusCodeBadRequest represents HTTP status code 400
	StatusCodeBadRequest = http.StatusBadRequest
	// StatusCodeInternalServerError represents HTTP status code 500
	StatusCodeInternalServerError = http.StatusInternalServerError
)

// AddSubscriptionHandler handles the endpoint for adding a new subscription.
func AddSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var newSubscription models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&newSubscription); err != nil {
		http.Error(w, "Invalid request body", StatusCodeBadRequest)
		return
	}

	userID := getUserIDFromRequest(r)
	if !isValidUserID(userID) {
		http.Error(w, "Invalid user ID", StatusCodeBadRequest)
		return
	}

	newSubscription.UserID = userID
	newSubscription.CreatedAt = time.Now()

	if _, err := db.SubscribeCollection.InsertOne(db.Ctx, newSubscription); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add subscription: %v", err), StatusCodeInternalServerError)
		return
	}

	w.WriteHeader(StatusCodeOK)
}

// GetSubscriptionsHandler handles the endpoint for fetching user subscriptions.
func GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromRequest(r)
	if !isValidUserID(userID) {
		http.Error(w, "Invalid user ID", StatusCodeBadRequest)
		return
	}

	cursor, err := db.SubscribeCollection.Find(db.Ctx, bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch subscriptions: %v", err), StatusCodeInternalServerError)
		return
	}
	defer cursor.Close(db.Ctx)

	var subscriptions []models.Subscription
	if err := cursor.All(db.Ctx, &subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode subscriptions: %v", err), StatusCodeInternalServerError)
		return
	}

	response, err := json.Marshal(subscriptions)
	if err != nil {
		http.Error(w, "Failed to serialize subscriptions", StatusCodeInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// getUserIDFromRequest retrieves the user ID from the request (customize as needed).
func getUserIDFromRequest(r *http.Request) string {
	// Example: Extract user ID from JWT token or session
	return "user_id_here"
}

// isValidUserID checks if the user ID is valid (customize as needed).
func isValidUserID(userID string) bool {
	// Add your validation logic here
	return userID != ""
}
