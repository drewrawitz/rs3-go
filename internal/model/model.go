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
