# Multiplayer Dot Game

A simple real-time multiplayer game where players control colored dots on a shared canvas.

## 🎮 Features

- Real-time multiplayer gameplay using WebSockets
- Each player gets a randomly colored dot
- Move your dot using arrow keys or WASD
- See other players' dots in real-time
- Automatic reconnection if connection is lost
- Boundary collision detection

## 🚀 Quick Start

### Option 1: GitHub Pages Client (Recommended for players)
- Visit the live client: `https://yourusername.github.io/yourrepo/`
- Enter the server URL provided by the host
- Start playing!

### Option 2: Local Development
1. Install Go dependencies:
   ```
   go mod tidy
   ```

2. Start the server:
   ```
   go run main.go cert.go
   ```

3. Open your browser and go to:
   ```
   http://localhost:3000
   ```

## 🖥️ For Server Hosts

### Running the Server
```bash
go run main.go cert.go
```

The server will start on:
- **HTTP**: `http://localhost:3000` (local development)
- **HTTPS**: `https://localhost:8443` (for GitHub Pages clients)

### Sharing with Network Players

Find your local IP address:
```bash
ipconfig | findstr "IPv4"    # Windows
ifconfig | grep "inet "      # Mac/Linux
```

Share these URLs with players:
- **HTTP**: `ws://YOUR_IP:3000/ws`
- **HTTPS**: `wss://YOUR_IP:8443/ws`

### Network Access Setup

For others on your network to connect:
1. Find your local IP address
2. Share the appropriate WebSocket URL
3. Configure Windows Firewall (if needed):
   ```powershell
   New-NetFirewallRule -DisplayName "Go Game Server" -Direction Inbound -Protocol TCP -LocalPort 3000,8443 -Action Allow
   ```

## 📁 Project Structure

```
├── docs/                 # GitHub Pages client
│   ├── index.html       # Web client that connects to any server
│   └── README.md        # GitHub Pages documentation
├── static/              # Local development client  
│   └── index.html       # Auto-connects to localhost
├── main.go              # Go WebSocket server
├── cert.go              # HTTPS certificate generation
├── go.mod               # Go module dependencies
└── README.md            # This file
```

## 🎯 Controls

- **Arrow Keys** or **WASD**: Move your dot
- The dot with a white border is yours
- Other players appear as colored dots without borders

## 🔧 GitHub Pages Setup

1. Push this repository to GitHub
2. Go to repository Settings → Pages
3. Select "Deploy from a branch"
4. Choose "main" branch and "/docs" folder
5. Your client will be available at `https://yourusername.github.io/reponame/`

## 🌐 Company Network Solutions

If you're on a company network without admin privileges:
- Try different ports (3000, 8443, 8000, etc.)
- Use the GitHub Pages client to avoid local firewall issues
- Host can use mobile hotspot as alternative network

Enjoy playing! 🎮
