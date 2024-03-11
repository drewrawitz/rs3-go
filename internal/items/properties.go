package items

import (
	"math/rand"
	"time"
)

// PickaxeProperties specific to pickaxes.

type PickaxeProperties struct {
	Penetration int32 `bson:"penetration" json:"penetration"`
	MinLevel    int32 `bson:"minLevel" json:"minLevel"`
	MinDamage   int32 `bson:"minDamage" json:"minDamage"`
	MaxDamage   int32 `bson:"maxDamage" json:"maxDamage"`
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func (p *PickaxeProperties) CalculateDamageRoll() int32 {
	// Ensure maxDamage is greater to avoid negative range
	if p.MaxDamage > p.MinDamage {
		// Calculate random damage within the range
		damageRange := int(p.MaxDamage - p.MinDamage + 1) // +1 to make it inclusive
		return p.MinDamage + int32(rnd.Intn(damageRange))
	}

	// Return minDamage if maxDamage is not greater (to handle unexpected cases)
	return p.MinDamage
}
