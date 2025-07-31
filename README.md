# TChat

`TChat` is a simple TCP-based chat application written in Go. It provides a basic client-server architecture for real-time messaging over a local network.

## Features
- TCP server that handles multiple clients
- User registration and message broadcasting
- Simple event-based protocol for communication
- Logging with different log levels

## Project Structure
```
go.mod
client/
    main.go         # Client entry point
    readers.go      # User input and connection helpers
internal/
    events/
        events.go   # Event types and serialization
    log/
        log.go      # Logging utility
server/
    args.go         # Command-line argument parsing
    chat.go         # Chat and user management
    main.go         # Server entry point
```

## Getting Started

### Prerequisites
- Go 1.24.5 or later

### Build
```
go build -o tchat-server ./server
go build -o tchat-client ./client
```

### Run the Server
```
./tchat-server <port>
```
Example:
```
./tchat-server 9000
```

### Run the Client
```
./tchat-client
```
You will be prompted for the server host, port, and your username. Once connection is stablished; then you can chat with other people connected people.

## How It Works
- The server listens for TCP connections and manages connected users.
- Clients connect, register a username, and can send/receive messages.
- Communication is handled via serialized event objects.