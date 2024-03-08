package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Name        string      `bson:"name"`
	Value       int32       `bson:"value"`
	Equipable   bool        `bson:"equipable"`
	Description string      `bson:"description"`
	Type        string      `bson:"type"`
	Properties  interface{} `bson:"properties"`
}

type PickaxeProperties struct {
	Penetration int32 `bson:"penetration"`
	MinLevel    int32 `bson:"minLevel"`
	MinDamage   int32 `bson:"minDamage"`
	MaxDamage   int32 `bson:"maxDamage"`
}

type HatchetProperties struct {
	Power int32 `bson:"power"`
}

var dbClient *mongo.Client

func unmarshalProperties(props bson.M, target interface{}) error {
	b, err := bson.Marshal(props)
	if err != nil {
		return err
	}
	return bson.Unmarshal(b, target)
}

func getItemByName(itemName string) (*Item, error) {
	col := dbClient.Database("rs3").Collection("items")

	var rawItem bson.M
	if err := col.FindOne(context.Background(), bson.M{"name": itemName}).Decode(&rawItem); err != nil {
		return nil, err
	}

	properties := rawItem["properties"].(bson.M)

	item := &Item{
		Name:        rawItem["name"].(string),
		Value:       rawItem["value"].(int32),
		Equipable:   rawItem["equipable"].(bool),
		Description: rawItem["description"].(string),
		Type:        rawItem["type"].(string),
	}

	switch item.Type {
	case "pickaxe":
		pickaxeProps := &PickaxeProperties{}
		if err := unmarshalProperties(properties, pickaxeProps); err != nil {
			return nil, err
		}
		item.Properties = pickaxeProps

	default:
		return nil, fmt.Errorf("unknown item type: %s", item.Type)
	}

	return item, nil
}

func main() {
	mongoURI := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error
	dbClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = dbClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Defer the disconnection
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			log.Fatalf("Disconnect error: %v", err)
		}
	}()

	item, err := getItemByName("Bronze Pickaxe2")

	if err != nil {
		log.Fatalf("Error getting item: %v", err)
	}

	if pickaxeProps, ok := item.Properties.(*PickaxeProperties); ok {
		fmt.Printf("Successfully retrieved and asserted pickaxe properties: %+v\n", pickaxeProps.MinDamage)
	} else {
		log.Fatalf("Failed to assert pickaxe properties")
	}

}
