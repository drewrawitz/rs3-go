package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rs3/internal/db"
	"rs3/internal/items"
	"rs3/internal/model"

	// "rs3/internal/model"
	"rs3/internal/skills/mining"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	// "go.mongodb.org/mongo-driver/mongo"
)

var activeMiningSessions = make(map[string]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// FIXME: In production, we should check the origin of requests:
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type MiningRequest struct {
	Action string `json:"action"`
	RockId string `json:"rockId"`
}

func main() {
	mongoURI := "mongodb://localhost:27017"
	client, err := db.Connect(mongoURI)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	r := gin.Default()

	r.GET("/ws/mine", func(c *gin.Context) {
		// Upgrade the HTTP server connection to the WebSocket protocol
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
			return
		}
		defer ws.Close()

		userID := "defekt7x"

		// Handle the WebSocket connection
		for {
			_, message, err := ws.ReadMessage()

			if err != nil {
				log.Println("Read error:", err)
				break
			}

			// Unmarshal the JSON data into the MiningRequest struct
			var req MiningRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			switch req.Action {
			case "startMining":

				if activeMiningSessions[userID] {
					log.Println("User is already mining")
					continue
				}

				rock, err := mining.GetRockById(c, client, req.RockId)
				if err != nil {
					// Send error message back via WebSocket
					continue
				}

				pickaxe, _ := items.GetItemById(c, client, "bronzePickaxe")

				playerInventory := &model.Inventory{
					Items: make(map[string]*model.InventoryItem),
				}

				playerSkillXP := map[string]int{
					"Mining": 0,
				}

				playerContext := model.PlayerContext{
					Inventory: playerInventory,
					SkillXP:   playerSkillXP,
				}

				rock.Mine(pickaxe, &playerContext, func(msg []byte) {
					if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("Write error:", err)
					}
				})

				activeMiningSessions[userID] = true
				defer func() { delete(activeMiningSessions, userID) }()

				// response := MiningResponse{
				// 	Status:  "success",
				// 	Message: "Mining started successfully",
				// 	RockId:  req.RockId, // Assuming req is your parsed request struct
				// }
				//
				// // Marshal the response struct to JSON
				// jsonResponse, err := json.Marshal(response)
				// if err != nil {
				// 	log.Println("JSON marshal error:", err)
				// 	continue
				// }
				//
				// if err := ws.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				// 	log.Println("Write error:", err)
				// 	break
				// }
			default:
				log.Println("Unrecognizable action", req.Action)
				continue
			}
		}
	})

	// router.GET("/items/:id", func(c *gin.Context) {
	// 	itemName := c.Param("id")
	// 	item, err := items.GetItemById(c, client, itemName)
	//
	// 	if err != nil {
	// 		// Check if the error is because the item was not found
	// 		if err == mongo.ErrNoDocuments {
	// 			c.JSON(http.StatusNotFound, gin.H{"error": "Item does not exist"})
	// 		} else {
	// 			// For other errors, return a 500 status code
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 		}
	//
	// 		return
	// 	}
	//
	// 	c.JSON(http.StatusOK, item)
	// })
	//
	// router.POST("swing", func(c *gin.Context) {
	// 	rock, rockErr := mining.GetRockById(c, client, "copperRock")
	//
	// 	playerInventory := &model.Inventory{
	// 		Items: make(map[string]*model.InventoryItem),
	// 	}
	//
	// 	playerSkillXP := map[string]int{
	// 		"Mining": 0,
	// 	}
	//
	// 	playerContext := model.PlayerContext{
	// 		Inventory: playerInventory,
	// 		SkillXP:   playerSkillXP,
	// 	}
	//
	// 	if rockErr != nil {
	// 		// Check if the error is because the item was not found
	// 		if err == mongo.ErrNoDocuments {
	// 			c.JSON(http.StatusNotFound, gin.H{"error": "Item does not exist"})
	// 		} else {
	// 			// For other errors, return a 500 status code
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 		}
	//
	// 		return
	// 	}
	//
	// 	pickaxe, pickaxeErr := items.GetItemById(c, client, "bronzePickaxe")
	//
	// 	if pickaxeErr != nil {
	// 		// Handle or log error when failing to get pickaxe
	// 		log.Printf("Error fetching pickaxe: %v", pickaxeErr)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pickaxe"})
	// 		return
	// 	}
	//
	// 	mineErr := rock.Mine(pickaxe, &playerContext)
	// 	if mineErr != nil {
	// 		// Log the error and/or return an appropriate response
	// 		log.Printf("Error mining rock: %v", mineErr)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": mineErr.Error()})
	// 		return
	// 	}
	//
	// 	c.JSON(http.StatusOK, rock)
	// })

	r.Run()
}

