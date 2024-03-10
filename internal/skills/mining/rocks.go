package mining

import (
	"context"
	"encoding/json"
	// "errors"
	"fmt"
	"log"
	"rs3/internal/items"
	"rs3/internal/model"
	"time"

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

const SWING_TIME = 2400 * time.Millisecond

func GetRockById(ctx context.Context, client *mongo.Client, rockId string) (*Rock, error) {
	col := client.Database("rs3").Collection("rocks")

	var rock Rock
	if err := col.FindOne(ctx, bson.M{"_id": rockId}).Decode(&rock); err != nil {
		return nil, err
	}

	return &rock, nil
}

func (r *Rock) Mine(pickaxe *model.Item, ctx *model.PlayerContext, sendMessage func(msg []byte)) {
	level := int32(3)

	// Check if a pickaxe is equipped
	if pickaxe == nil {
		progressMessage := MiningResponse{
			Status:  "error",
			Message: "You don't have a pickaxe equipped",
			RockId:  r.ID,
		}
		progressJSON, err := json.Marshal(progressMessage)

		if err != nil {
			log.Printf("Error marshalling progress message: %v", err)
			return
		}
		sendMessage(progressJSON)
		return
	}

	pickaxeProps, isValidPickaxe := pickaxe.Properties.(items.PickaxeProperties)

	if !isValidPickaxe {

		progressMessage := MiningResponse{
			Status:  "error",
			Message: "Item used to mine the rock is not a pickaxe",
			RockId:  r.ID,
		}
		progressJSON, _ := json.Marshal(progressMessage)
		sendMessage(progressJSON)
		return
	}

	// Make sure the user has the appropriate mining level to mine this rock
	if level < r.MinLevel {
		progressMessage := MiningResponse{
			Status:  "error",
			Message: fmt.Sprintf("You must be level %d Mining to mine this rock.", r.MinLevel),
			RockId:  r.ID,
		}
		progressJSON, _ := json.Marshal(progressMessage)
		sendMessage(progressJSON)
		return
	}

	// Make sure the user has the appropriate mining level to use this pickaxe
	if level < pickaxeProps.MinLevel {
		progressMessage := MiningResponse{
			Status:  "error",
			Message: fmt.Sprintf("You must be level %d Mining to use this pickaxe.", pickaxeProps.MinLevel),
			RockId:  r.ID,
		}
		progressJSON, _ := json.Marshal(progressMessage)
		sendMessage(progressJSON)
		return
	}

	progressMessage := MiningResponse{
		Status:  "swinging",
		Message: fmt.Sprintf("You swing your %s at the %s", pickaxe.Name, r.Name),
		RockId:  r.ID,
	}
	progressJSON, _ := json.Marshal(progressMessage)
	sendMessage(progressJSON)

	ma := NewMiningActivity(level)

	for {
		health := r.Durability

		for health > 0 {
			time.Sleep(SWING_TIME)

			netHardness := int32(pickaxeProps.Penetration) - int32(r.Hardness)
			damageRoll := pickaxeProps.CalculateDamageRoll()
			damage := int32(level) + damageRoll + netHardness

			if damage >= int32(health) {
				progressMessage := MiningResponse{
					Status:  "swinging",
					Message: fmt.Sprintf("You mine the rock and deal %v damage.", health),
					Damage:  health,
					RockId:  r.ID,
				}
				progressJSON, _ := json.Marshal(progressMessage)
				sendMessage(progressJSON)

				// Process drops
				drops := ma.processDrops(r, ctx)

				a := MiningResponse{
					Status:  "swinging",
					Message: "Loot has been dropped",
					Drops:   drops,
				}
				progressJSON2, _ := json.Marshal(a)
				sendMessage(progressJSON2)

				break // Break the inner loop, rock is considered depleted
			}

			health -= int32(damage)

			progress := (float32(r.Durability-health) / float32(r.Durability)) * 100

			progressMessage := MiningResponse{
				Status:   "swinging",
				Message:  fmt.Sprintf("You mine the rock and deal %v damage.", damage),
				Progress: progress,
				Damage:   damage,
				RockId:   r.ID,
			}
			progressJSON, _ := json.Marshal(progressMessage)
			sendMessage(progressJSON)
		}

		fmt.Println("The rock resets for another round of mining.")
	}
}
