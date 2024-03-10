package mining

import (
	"rs3/internal/model"
)

type Rock struct {
	ID             string          `bson:"_id" json:"id"`
	Name           string          `bson:"name" json:"name"`
	Durability     int32           `bson:"durability" json:"durability"`
	Hardness       int32           `bson:"hardness" json:"hardness"`
	MinLevel       int32           `bson:"minLevel" json:"minLevel"`
	Multiplier     float64         `bson:"multiplier" json:"multiplier"`
	PrimaryDrops   model.DropTable `bson:"primaryDrops" json:"primaryDrops"`
	SecondaryDrops model.DropTable `bson:"secondaryDrops" json:"secondaryDrops"`
}
