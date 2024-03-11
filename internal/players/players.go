package players

import (
	"context"
	"fmt"
	"rs3/internal/items"
	"rs3/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPlayerByUsername(ctx context.Context, client *mongo.Client, username string) (*model.Player, error) {
	col := client.Database("rs3").Collection("players")

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "username", Value: username}}}},
		bson.D{{Key: "$unwind", Value: "$inventory"}},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "items"},
				{Key: "localField", Value: "inventory.itemId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "inventory.item"},
			}},
		},
		bson.D{{Key: "$unwind", Value: "$inventory.item"}},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "username", Value: bson.D{{Key: "$first", Value: "$username"}}},
				{Key: "inventory", Value: bson.D{{Key: "$push", Value: "$inventory"}}},
			}},
		},
	}

	var player model.Player
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no player found with username: %s", username)
	}

	return &player, nil
}

func FindHighestLevelPickaxe(player *model.Player) *model.InventoryItem {
	if player == nil {
		return nil
	}

	var highestPickaxe *model.InventoryItem
	highestLevel := int32(-1)

	for _, inventoryItem := range player.Inventory {
		if inventoryItem.Item.Type == "pickaxe" {

			if err := items.UnmarshalItemProperties(&inventoryItem.Item); err != nil {
				fmt.Println("Error unmarshalling item properties:", err)
				continue
			}

			if pickaxeProps, ok := inventoryItem.Item.Properties.(items.PickaxeProperties); ok {
				if pickaxeProps.MinLevel > highestLevel {
					highestLevel = pickaxeProps.MinLevel
					highestPickaxe = &inventoryItem
				}
			}
		}
	}

	return highestPickaxe
}
