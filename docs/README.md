# Multiplayer Dot Game - GitHub Pages Client

This is the web client for the Multiplayer Dot Game that can connect to any game server.

## 🎮 How to Play

1. **Host runs the server**: Someone needs to run the Go server on their computer
2. **Get the server URL**: The host shares their WebSocket URL (example: `ws://192.168.1.35:3000/ws`)
3. **Connect**: Enter the server URL and click "Connect"
4. **Play**: Move your dot with arrow keys or WASD keys

## 🔗 Live Demo

Visit the live client at: `https://yourusername.github.io/yourrepo/`

## 🖥️ For Server Hosts

To run the game server:

1. Clone this repository
2. Install Go (if not installed)
3. Run the server:
   ```bash
   cd path/to/repo
   go run main.go cert.go
   ```
4. Share your server URL with players:
   - HTTP: `ws://YOUR_IP:3000/ws`
   - HTTPS: `wss://YOUR_IP:8443/ws`

## 📁 Repository Structure

```
├── docs/                 # GitHub Pages client
│   └── index.html       # Web client that connects to any server
├── static/              # Local development client
│   └── index.html       # Auto-connects to localhost
├── main.go              # Go WebSocket server
├── cert.go              # HTTPS certificate generation
└── README.md            # Main project documentation
```

## 🚀 GitHub Pages Setup

1. Push this repository to GitHub
2. Go to repository Settings → Pages
3. Select "Deploy from a branch"
4. Choose "main" branch and "/docs" folder
5. Your client will be available at `https://yourusername.github.io/reponame/`

## 🔧 Network Requirements

- All players must be able to reach the host's computer
- Host may need to configure firewall/router for external access
- For company networks, try different ports (3000, 8443, etc.)

Enjoy the game! 🎯
