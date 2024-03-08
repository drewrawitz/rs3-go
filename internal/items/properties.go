package items

// PickaxeProperties specific to pickaxes.
type PickaxeProperties struct {
	Penetration int32 `bson:"penetration"`
	MinLevel    int32 `bson:"minLevel"`
	MinDamage   int32 `bson:"minDamage"`
	MaxDamage   int32 `bson:"maxDamage"`
}
