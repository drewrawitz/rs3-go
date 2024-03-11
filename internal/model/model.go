package model

// Item represents a general game item.
type Item struct {
	ID          string      `bson:"_id" json:"id"`
	Name        string      `bson:"name" json:"name"`
	Value       int         `bson:"value" json:"value"`
	Equipable   bool        `bson:"equipable" json:"equipable"`
	Description string      `bson:"description" json:"description"`
	Type        string      `bson:"type" json:"type"`
	Properties  interface{} `bson:"properties" json:"properties"`
}

type InventoryItem struct {
	ItemID string `bson:"itemId" json:"itemId"`
	Qty    int    `bson:"qty" json:"qty"`
	Item   Item   `bson:"item" json:"item"`
}

type LootItem struct {
	ItemID    string  `bson:"itemId" json:"itemId"`
	MinQty    int     `bson:"minQty" json:"minQty"`
	MaxQty    int     `bson:"maxQty" json:"maxQty"`
	Chance    float64 `bson:"chance" json:"chance"`
	Condition string  `bson:"condition" json:"condition"`
}

type DropTable []LootItem

type Player struct {
	ID        string          `bson:"_id" json:"id"`
	Username  string          `bson:"username" json:"username"`
	Inventory []InventoryItem `bson:"inventory" json:"inventory"`
}
