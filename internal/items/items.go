package items

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rs3/internal/model"
)

func unmarshalProperties(props bson.M, target interface{}) error {
	b, err := bson.Marshal(props)
	if err != nil {
		return err
	}
	return bson.Unmarshal(b, target)
}

// getItemByName fetches an item by its name from the MongoDB database.
func GetItemByName(client *mongo.Client, itemName string) (*model.Item, error) {
	col := client.Database("rs3").Collection("items")

	var rawItem bson.M
	if err := col.FindOne(context.Background(), bson.M{"name": itemName}).Decode(&rawItem); err != nil {
		return nil, err
	}

	item := &model.Item{
		Name:        rawItem["name"].(string),
		Value:       rawItem["value"].(int32),
		Equipable:   rawItem["equipable"].(bool),
		Description: rawItem["description"].(string),
		Type:        rawItem["type"].(string),
	}

	// Handle properties based on item.Type
	switch item.Type {
	case "pickaxe":
		pickaxeProps := &PickaxeProperties{}
		if err := unmarshalProperties(rawItem["properties"].(bson.M), pickaxeProps); err != nil {
			return nil, err
		}
		item.Properties = pickaxeProps
	default:
		return nil, fmt.Errorf("unknown item type: %s", item.Type)
	}

	return item, nil
}
