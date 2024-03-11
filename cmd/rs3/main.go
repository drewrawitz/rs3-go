package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rs3/internal/db"
	"rs3/internal/players"
	"rs3/internal/skills/mining"
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
				rock, err := mining.GetRockById(c, client, req.RockId)
				if err != nil {
					// Send error message back via WebSocket
					continue
				}

				player, _ := players.GetPlayerByUsername(c, client, "DevEx7")

				miningSession := &mining.MiningSession{
					Rock:       rock,
					Player:     player,
					SendMsg:    func(msg []byte) { ws.WriteMessage(websocket.TextMessage, msg) },
					StopSignal: make(chan bool),
				}

				miningSession.Start() // Start the mining session

				// rock.Mine(pickaxe, &playerContext, func(msg []byte) {
				// 	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				// 		log.Println("Write error:", err)
				// 	}
				// })
				//
				// activeMiningSessions[userID] = true
				// defer func() { delete(activeMiningSessions, userID) }()

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

	r.Run()
}
