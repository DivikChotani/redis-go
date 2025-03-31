# Redis-Go (Minimal TCP Server)

This project is a minimal TCP server written in Go that mimics the basic structure of a Redis server. It listens on port `6379` (the default Redis port), accepts a client connection, logs incoming requests, and responds with a simple `+OK`.

## Current Features (will be improved upon)

- Starts a TCP server on `localhost:6379`
- Accepts a single incoming client connection
- Reads client input using a buffer buffer and logs the request
- Sends back a hardcoded Redis-style response: `+OK\r\n`



