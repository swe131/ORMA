// mongo.go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Subscription represents the structure of your MongoDB document
type SubscribeCollection struct {
	// Define your subscription fields here
	// For example:
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"user_id"`
	Name   string             `bson:"name"`
	Amount float64            `bson:"amount"`
	// Add other fields as needed
}

var (
	Client        *mongo.Client
	ctx           context.Context
	Subscriptions *mongo.Collection // Move this line here
)

// ConnectDB establishes a connection to the MongoDB server
func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.TODO()
}

// InitMongoDB initializes the MongoDB connection and collection
func InitMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Connected to MongoDB!")

	// Set the Subscriptions collection
	Subscriptions = client.Database("orma").Collection("subscriptions")

	log.Println("Subscriptions collection initialized!")

	return nil
}

// GetAllSubscriptions retrieves all subscriptions from the database
func GetAllSubscriptions() ([]SubscribeCollection, error) {
	var subscriptions []SubscribeCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := Subscriptions.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Other functions...
