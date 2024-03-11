package mining

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MiningResponse struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	RockId   string      `json:"rockId,omitempty"`
	Progress float32     `json:"progress,omitempty"`
	Damage   int32       `json:"damage,omitempty"`
	Drops    interface{} `json:"drops,omitempty"`
}

func GetRockById(ctx context.Context, client *mongo.Client, rockId string) (*Rock, error) {
	col := client.Database("rs3").Collection("rocks")

	var rock Rock
	if err := col.FindOne(ctx, bson.M{"_id": rockId}).Decode(&rock); err != nil {
		return nil, err
	}

	return &rock, nil
}
