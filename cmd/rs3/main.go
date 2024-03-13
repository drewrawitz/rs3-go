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

var activeMiningSessions = make(map[string]*mining.MiningSession)

func stopAndCleanupSession(playerID string) {
	if session, exists := activeMiningSessions[playerID]; exists {
		// Send a signal to stop the specific player's mining session
		session.StopSignal <- true

		// Remove the session from the active sessions map
		delete(activeMiningSessions, playerID)
	} else {
		log.Printf("No active mining session found for player ID %s to stop", playerID)
	}
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

		player, _ := players.GetPlayerByUsername(c, client, "DevEx7")

		// Handle the WebSocket connection
		for {
			_, message, err := ws.ReadMessage()

			if err != nil {
				log.Println("WebSocket disconnected:", err)
				stopAndCleanupSession(player.ID)
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

				miningSession := &mining.MiningSession{
					Rock:       rock,
					Client:     client,
					Context:    c,
					Player:     player,
					SendMsg:    func(msg []byte) { ws.WriteMessage(websocket.TextMessage, msg) },
					StopSignal: make(chan bool),
				}

				activeMiningSessions[player.ID] = miningSession
				miningSession.Start()

			case "stopMining":
				stopAndCleanupSession(player.ID)

			default:
				log.Println("Unrecognizable action", req.Action)
				continue
			}
		}
	})

	r.Run()
}
