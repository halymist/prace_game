package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Player struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color string  `json:"color"`
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type MoveData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type GameState struct {
	Players map[string]*Player `json:"players"`
}

var (
	players = make(map[string]*Player)
	mutex   = sync.RWMutex{}
	colors  = []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00", "#FF00FF", "#00FFFF", "#FFA500", "#800080"}
)

func generatePlayerID() string {
	return fmt.Sprintf("player_%d", time.Now().UnixNano())
}

func getRandomColor() string {
	return colors[rand.Intn(len(colors))]
}

func broadcastGameState() {
	mutex.RLock()
	defer mutex.RUnlock()

	gameState := GameState{Players: make(map[string]*Player)}
	for id, player := range players {
		gameState.Players[id] = &Player{
			ID:    player.ID,
			X:     player.X,
			Y:     player.Y,
			Color: player.Color,
		}
	}

	message := Message{
		Type: "gameState",
		Data: gameState,
	}

	data, _ := json.Marshal(message)

	for _, player := range players {
		player.Mutex.Lock()
		if player.Conn != nil {
			err := player.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error sending message to player %s: %v", player.ID, err)
			}
		}
		player.Mutex.Unlock()
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	playerID := generatePlayerID()
	player := &Player{
		ID:    playerID,
		X:     float64(rand.Intn(750) + 25), // Random position within canvas bounds
		Y:     float64(rand.Intn(550) + 25),
		Color: getRandomColor(),
		Conn:  conn,
	}

	mutex.Lock()
	players[playerID] = player
	mutex.Unlock()

	log.Printf("Player %s connected", playerID)

	// Send player their own ID first
	welcomeMsg := Message{
		Type: "welcome",
		Data: map[string]string{"playerId": playerID},
	}
	welcomeData, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, welcomeData)

	// Send initial game state
	broadcastGameState()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Player %s disconnected: %v", playerID, err)
			break
		}

		switch msg.Type {
		case "move":
			var moveData MoveData
			if data, ok := msg.Data.(map[string]interface{}); ok {
				if x, ok := data["x"].(float64); ok {
					moveData.X = x
				}
				if y, ok := data["y"].(float64); ok {
					moveData.Y = y
				}
			}

			mutex.Lock()
			if player, exists := players[playerID]; exists {
				// Boundary checking
				if moveData.X >= 10 && moveData.X <= 790 {
					player.X = moveData.X
				}
				if moveData.Y >= 10 && moveData.Y <= 590 {
					player.Y = moveData.Y
				}
			}
			mutex.Unlock()

			broadcastGameState()
		}
	}

	// Remove player when disconnected
	mutex.Lock()
	delete(players, playerID)
	mutex.Unlock()

	broadcastGameState()
	log.Printf("Player %s removed", playerID)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate self-signed certificate if it doesn't exist
	if _, err := os.Stat("server.crt"); os.IsNotExist(err) {
		log.Println("Generating self-signed certificate...")
		if err := generateSelfSignedCert(); err != nil {
			log.Printf("Failed to generate certificate: %v", err)
		}
	}

	// Serve docs files (works for both local testing and GitHub Pages)
	http.Handle("/", http.FileServer(http.Dir("./docs/")))

	// WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	// Start HTTPS server
	go func() {
		fmt.Println("HTTPS Server starting on :8443")
		fmt.Println("Open https://localhost:8443 in your browser")
		fmt.Println("For network access: https://192.168.1.35:8443")
		if err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil); err != nil {
			log.Printf("HTTPS server failed: %v", err)
		}
	}()

	// Start HTTP server as fallback
	fmt.Println("HTTP Server starting on :3000")
	fmt.Println("Open http://localhost:3000 in your browser")
	fmt.Println("For network access: http://192.168.1.35:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
