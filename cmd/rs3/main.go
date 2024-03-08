package main

import (
	"log"
	"rs3/internal/db"
	"rs3/internal/items"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	client, err := db.Connect(mongoURI)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := db.Disconnect(client); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	item, err := items.GetItemByName(client, "Bronze Pickaxe")
	if err != nil {
		log.Fatalf("Error getting item: %v", err)
	}

	if pickaxeProps, ok := item.Properties.(*items.PickaxeProperties); ok {
		log.Printf("Successfully retrieved and asserted pickaxe properties: %+v\n", pickaxeProps)
	} else {
		log.Fatalf("Failed to assert pickaxe properties")
	}
}
