package mining

import (
	// "log"
	// "math/rand"
	// dropUtils "rs3/internal/game/drops"
	"log"
	"rs3/internal/model"
	"time"
)

type SessionStats struct {
	TotalOres          int32
	TotalRocksDepleted int32
	MiningDuration     time.Duration
}

type MiningActivity struct {
	CurrentRock  *Rock
	MiningLevel  int32
	SessionStats SessionStats
}

func NewMiningActivity(playerLevel int32) *MiningActivity {
	return &MiningActivity{
		CurrentRock:  nil,
		MiningLevel:  playerLevel,
		SessionStats: SessionStats{TotalOres: 0, TotalRocksDepleted: 0, MiningDuration: 0},
	}
}

func (ma *MiningActivity) processDrops(rock *Rock, ctx *model.PlayerContext) model.DropTable {
	log.Print(ctx)
	return rock.PrimaryDrops
	// processConditionalDrops := func(drops model.DropTable) {
	// 	for _, drop := range drops {
	// 		shouldDrop := rand.Float64() < drop.Chance
	//
	// 		if !shouldDrop {
	// 			continue
	// 		}
	//
	// 		quantity := dropUtils.CalculateQuantity(drop.MinQty, drop.MaxQty)
	//
	// 		log.Printf("Add item to inventory: %s (%v)", drop.ItemID, quantity)
	// 	}
	// }
	//
	// processConditionalDrops(rock.PrimaryDrops)
	// processConditionalDrops(rock.SecondaryDrops)
}
