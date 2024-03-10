package model

// Item represents a general game item.
type Item struct {
	Name        string      `bson:"name"`
	Value       int32       `bson:"value"`
	Equipable   bool        `bson:"equipable"`
	Description string      `bson:"description"`
	Type        string      `bson:"type"`
	Properties  interface{} `bson:"properties"`
}

type InventoryItem struct {
	Item     *Item
	Quantity int
}

type Inventory struct {
	Items map[string]*InventoryItem
}

type LootItem struct {
	ItemID    string  `bson:"itemId" json:"itemId"`
	MinQty    int     `bson:"minQty" json:"minQty"`
	MaxQty    int     `bson:"maxQty" json:"maxQty"`
	Chance    float64 `bson:"chance" json:"chance"`
	Condition string  `bson:"condition" json:"condition"`
}

type DropTable []LootItem

type PlayerContext struct {
	Inventory *Inventory
	SkillXP   map[string]int
}
