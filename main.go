package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
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
	HP    int     `json:"hp"`
	MaxHP int     `json:"maxHP"`
	Angle float64 `json:"angle"` // Cannon direction in radians
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

type ShootData struct {
	Angle float64 `json:"angle"`
}

type Projectile struct {
	ID      string  `json:"id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	VelX    float64 `json:"velX"`
	VelY    float64 `json:"velY"`
	OwnerID string  `json:"ownerId"`
	Color   string  `json:"color"`
}

type GameState struct {
	Players     map[string]*Player     `json:"players"`
	Projectiles map[string]*Projectile `json:"projectiles"`
}

var (
	players     = make(map[string]*Player)
	projectiles = make(map[string]*Projectile)
	mutex       = sync.RWMutex{}
	colors      = []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00", "#FF00FF", "#00FFFF", "#FFA500", "#800080"}
)

func generatePlayerID() string {
	return fmt.Sprintf("player_%d", time.Now().UnixNano())
}

func generateProjectileID() string {
	return fmt.Sprintf("projectile_%d", time.Now().UnixNano())
}

func getRandomColor() string {
	return colors[rand.Intn(len(colors))]
}

func broadcastGameState() {
	mutex.RLock()
	defer mutex.RUnlock()

	gameState := GameState{
		Players:     make(map[string]*Player),
		Projectiles: make(map[string]*Projectile),
	}

	for id, player := range players {
		gameState.Players[id] = &Player{
			ID:    player.ID,
			X:     player.X,
			Y:     player.Y,
			Color: player.Color,
			HP:    player.HP,
			MaxHP: player.MaxHP,
			Angle: player.Angle,
		}
	}

	for id, projectile := range projectiles {
		gameState.Projectiles[id] = &Projectile{
			ID:      projectile.ID,
			X:       projectile.X,
			Y:       projectile.Y,
			VelX:    projectile.VelX,
			VelY:    projectile.VelY,
			OwnerID: projectile.OwnerID,
			Color:   projectile.Color,
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
		HP:    5,
		MaxHP: 5,
		Angle: 0,
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
		case "shoot":
			var shootData ShootData
			if data, ok := msg.Data.(map[string]interface{}); ok {
				if angle, ok := data["angle"].(float64); ok {
					shootData.Angle = angle
				}
			}

			mutex.Lock()
			if shooter, exists := players[playerID]; exists && shooter.HP > 0 {
				// Update player's cannon angle
				shooter.Angle = shootData.Angle

				// Create projectile
				projectileID := generateProjectileID()
				speed := 200.0
				projectile := &Projectile{
					ID:      projectileID,
					X:       shooter.X,
					Y:       shooter.Y,
					VelX:    speed * math.Cos(shootData.Angle),
					VelY:    speed * math.Sin(shootData.Angle),
					OwnerID: playerID,
					Color:   shooter.Color,
				}
				projectiles[projectileID] = projectile
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

func updateProjectiles() {
	mutex.Lock()
	defer mutex.Unlock()

	toRemove := []string{}

	for id, projectile := range projectiles {
		// Move projectile
		projectile.X += projectile.VelX * 0.016 // ~60fps
		projectile.Y += projectile.VelY * 0.016

		// Check bounds
		if projectile.X < 0 || projectile.X > 800 || projectile.Y < 0 || projectile.Y > 600 {
			toRemove = append(toRemove, id)
			continue
		}

		// Check collision with players
		for _, player := range players {
			if player.ID == projectile.OwnerID || player.HP <= 0 {
				continue
			}

			// Simple distance collision
			dx := projectile.X - player.X
			dy := projectile.Y - player.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance < 20 { // Hit!
				player.HP--
				if player.HP <= 0 {
					log.Printf("Player %s was eliminated by %s!", player.ID, projectile.OwnerID)
				}
				toRemove = append(toRemove, id)
				break
			}
		}
	}

	// Remove projectiles
	for _, id := range toRemove {
		delete(projectiles, id)
	}
}

func gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60fps
	defer ticker.Stop()

	for range ticker.C {
		updateProjectiles()
		broadcastGameState()
	}
}

func main() {
	// Serve docs files
	http.Handle("/", http.FileServer(http.Dir("./docs/")))

	// WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	// Start game loop
	go gameLoop()

	// Start HTTP server
	fmt.Println("HTTP Server starting on :3000")
	fmt.Println("Open http://localhost:3000 in your browser")
	fmt.Println("For network access: http://10.0.1.14:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
