# Multiplayer Dot Game

A simple real-time multiplayer game where players control colored dots on a shared canvas.

## Features

- Real-time multiplayer gameplay using WebSockets
- Each player gets a randomly colored dot
- Move your dot using arrow keys or WASD
- See other players' dots in real-time
- Automatic reconnection if connection is lost
- Boundary collision detection

## How to Run

1. Install Go dependencies:
   ```
   go mod tidy
   ```

2. Start the server:
   ```
   go run main.go
   ```

3. Open your browser and go to:
   ```
   http://localhost:8080
   ```

4. Share the URL with others on your network:
   ```
   http://[YOUR_IP]:8080
   ```

## Controls

- **Arrow Keys** or **WASD**: Move your dot
- The dot with a white border is yours
- Other players appear as colored dots without borders

## Network Access

To allow others on your network to connect:
1. Find your local IP address
2. Share `http://[YOUR_IP]:8080` with other players
3. Make sure Windows Firewall allows the connection on port 8080

Enjoy playing!
