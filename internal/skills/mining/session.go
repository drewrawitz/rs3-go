package mining

import (
	"encoding/json"
	"fmt"
	"log"
	"rs3/internal/items"
	"rs3/internal/model"
	"rs3/internal/players"
	"time"
)

const SWING_TIME = 2400 * time.Millisecond

type MiningSession struct {
	Rock       *Rock
	Player     *model.Player
	SendMsg    func(msg []byte)
	StopSignal chan bool // Used to signal when the mining session should stop
}

func (ms *MiningSession) Start() {
	level := int32(30)
	pickaxe := players.FindHighestLevelPickaxe(ms.Player)
	initialRockHealth := ms.Rock.Durability

	go func() {
		ticker := time.NewTicker(SWING_TIME)
		defer ticker.Stop()

		for range ticker.C {
			if pickaxe == nil {
				ms.sendProgressMessage("error", "You don't have a pickaxe equipped", 0, 0)
				return
			}

			pickaxeProps, isValidPickaxe := pickaxe.Item.Properties.(items.PickaxeProperties)
			if !isValidPickaxe {
				ms.sendProgressMessage("error", "Item used to mine the rock is not a pickaxe", 0, 0)
				return
			}

			netHardness := int32(pickaxeProps.Penetration) - int32(ms.Rock.Hardness)
			damageRoll := pickaxeProps.CalculateDamageRoll()
			damage := int32(level) + damageRoll + netHardness

			// Ensure the last hit is the remaining Durability (don't go below 0)
			if (ms.Rock.Durability - damage) < 0 {
				damage = ms.Rock.Durability
			}

			ms.Rock.Durability -= damage
			progress := (float32(initialRockHealth-ms.Rock.Durability) / float32(initialRockHealth)) * 100

			ms.sendProgressMessage("swinging", fmt.Sprintf("You mine the rock and deal %v damage.", damage), progress, damage)

			if ms.Rock.Durability <= 0 {
				ms.sendProgressMessage("depleted", "The rock has been depleted.", 0, 0)
				ms.Rock.Durability = initialRockHealth
				continue
			}
		}
	}()
}

// Helper method to encapsulate message sending
func (ms *MiningSession) sendProgressMessage(status, message string, progress float32, damage int32) {
	progressMessage := MiningResponse{
		Status:   status,
		Message:  message,
		RockId:   ms.Rock.ID,
		Progress: progress,
		Damage:   int32(damage),
	}
	progressJSON, err := json.Marshal(progressMessage)
	if err != nil {
		log.Printf("Error marshaling progress message: %v\n", err)
		return
	}
	ms.SendMsg(progressJSON)
}
