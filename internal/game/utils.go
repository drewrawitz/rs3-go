package game

import (
	"go.mongodb.org/mongo-driver/bson"
)

func UnmarshalProperties(props bson.M, target interface{}) error {
	b, err := bson.Marshal(props)
	if err != nil {
		return err
	}
	return bson.Unmarshal(b, target)
}
