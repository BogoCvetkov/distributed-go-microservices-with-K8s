package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func New(mongo *mongo.Client) *Models {
	client = mongo

	return &Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Data      string    `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type LogEntryUpdate struct {
	Name string `json:"name" bson:"name,omitempty"`
	Data string `json:"data" bson:"data,omitempty"`
}

func (l *LogEntry) Insert(data *LogEntry) error {

	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add creation date
	data.CreatedAt = time.Now()

	_, err := collection.InsertOne(ctx, data)

	if err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) GetAll() (*[]LogEntry, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})

	var results []LogEntry

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var doc LogEntry

		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}

		results = append(results, doc)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &results, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var document LogEntry

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid mongo id")
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&document)
	if err == mongo.ErrNoDocuments {
		fmt.Println(err)
		// Do something when no record was found
		return nil, errors.New("record not found")
	} else if err != nil {
		return nil, err
	}

	return &document, nil
}

func (l *LogEntry) Update(id string, data *LogEntryUpdate) (*mongo.UpdateResult, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid mongo id")
	}

	update := bson.M{
		"$set": bson.M{
			"name":       data.Name,
			"data":       data.Data,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateByID(ctx, objID, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetCollection() *mongo.Collection {
	collection := client.Database("go_microservices").Collection("logs")
	return collection
}
