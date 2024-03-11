package drops

import (
	"math/rand"
)

// Helper function to calculate drop quantity based on min and max values
func CalculateQuantity(minQty, maxQty int) int {
	if minQty == maxQty {
		return minQty
	}

	return minQty + rand.Intn(maxQty-minQty+1)
}
