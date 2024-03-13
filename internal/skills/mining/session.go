package mining

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rs3/internal/items"
	"rs3/internal/model"
	"rs3/internal/players"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const SWING_TIME = 2400 * time.Millisecond

type MiningSession struct {
	Rock       *Rock
	Player     *model.Player
	Client     *mongo.Client
	Context    context.Context
	SendMsg    func(msg []byte)
	StopSignal chan bool // Used to signal when the mining session should stop
}

func (ms *MiningSession) Start() {
	var accumulatedXP float64 = 0

	level := int32(1)
	pickaxe := players.FindHighestLevelPickaxe(ms.Player)
	initialRockHealth := ms.Rock.Durability

	go func() {
		ticker := time.NewTicker(SWING_TIME)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if pickaxe == nil {
					ms.sendProgressMessage("error", "You don't have a pickaxe equipped", 0, 0, 0)
					return
				}

				pickaxeProps, isValidPickaxe := pickaxe.Item.Properties.(items.PickaxeProperties)
				if !isValidPickaxe {
					ms.sendProgressMessage("error", "Item used to mine the rock is not a pickaxe", 0, 0, 0)
					return
				}

				netHardness := int32(pickaxeProps.Penetration) - int32(ms.Rock.Hardness)
				damageRoll := pickaxeProps.CalculateDamageRoll()
				damage := int32(level) + damageRoll + netHardness

				// Ensure the last hit is the remaining Durability (don't go below 0)
				if (ms.Rock.Durability - damage) < 0 {
					damage = ms.Rock.Durability
				}

				// Calculate the experience gained
				xp := ms.Rock.calculateExperience(damage)
				accumulatedXP += xp

				// Damage the rock
				ms.Rock.Durability -= damage
				progress := (float32(initialRockHealth-ms.Rock.Durability) / float32(initialRockHealth)) * 100

				ms.sendProgressMessage("swinging", fmt.Sprintf("You mine the rock and deal %v damage.", damage), progress, damage, xp)

				if ms.Rock.Durability <= 0 {
					// Update the player's skill XP in MongoDB
					ms.updatePlayerSkillXP(accumulatedXP)
					accumulatedXP = 0

					// Reset the rock health
					ms.Rock.Durability = initialRockHealth
					continue
				}

			case <-ms.StopSignal:
				// Clean-up logic here
				log.Println("Stopping mining session for player", ms.Player.ID)

				// Update the player's skill XP in MongoDB
				ms.updatePlayerSkillXP(accumulatedXP)
				accumulatedXP = 0

				// Stop the Ticker
				ticker.Stop()
				return
			}
		}
	}()
}

// Helper method to encapsulate message sending
func (ms *MiningSession) sendProgressMessage(status, message string, progress float32, damage int32, xp float64) {
	progressMessage := MiningResponse{
		Status:   status,
		Message:  message,
		RockId:   ms.Rock.ID,
		Progress: progress,
		Damage:   int32(damage),
		Xp:       xp,
	}
	progressJSON, err := json.Marshal(progressMessage)
	if err != nil {
		log.Printf("Error marshaling progress message: %v\n", err)
		return
	}
	ms.SendMsg(progressJSON)
}

func (ms *MiningSession) updatePlayerSkillXP(xp float64) {
	playerCollection := ms.Client.Database("rs3").Collection("players")
	id, _ := primitive.ObjectIDFromHex(ms.Player.ID)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"skills.mining": xp}}

	_, err := playerCollection.UpdateOne(ms.Context, filter, update)
	if err != nil {
		log.Printf("Failed to update player skill XP: %v\n", err)
	}
}
