package mining

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const GLOBAL_MINING_MODIFIER = 0.4

type MiningResponse struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	RockId   string      `json:"rockId,omitempty"`
	Progress float32     `json:"progress,omitempty"`
	Damage   int32       `json:"damage,omitempty"`
	Xp       float64     `json:"xp,omitempty"`
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

func (r *Rock) calculateExperience(damage int32) float64 {
	return math.Ceil(float64(damage)*r.Multiplier*GLOBAL_MINING_MODIFIER*100) / 100
}
