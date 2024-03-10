package items

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rs3/internal/model"
)

func GetItemById(ctx context.Context, client *mongo.Client, itemId string) (*model.Item, error) {
	col := client.Database("rs3").Collection("items")

	var item model.Item
	if err := col.FindOne(ctx, bson.M{"_id": itemId}).Decode(&item); err != nil {
		return nil, err
	}

	switch item.Type {
	case "pickaxe":

		propertiesBSON, err := bson.Marshal(item.Properties)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal properties: %w", err)
		}

		var pickaxeProps PickaxeProperties
		if err := bson.Unmarshal(propertiesBSON, &pickaxeProps); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pickaxe properties: %w", err)
		}

		item.Properties = pickaxeProps

	default:
		return nil, fmt.Errorf("unknown item type: %s", item.Type)
	}

	return &item, nil
}
