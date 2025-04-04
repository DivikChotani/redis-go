# Redis-Go (Minimal TCP Server)

This project is a minimal TCP server written in Go that mimics the basic structure of a Redis server. It listens on port `6379` (the default Redis port), accepts a client connection, logs incoming requests, and responds with a simple `+OK`.

## Current Features (will be improved upon)

- Starts a TCP server on `localhost:6379`
- Accepts a single incoming client connection
- Reads client input using a buffer buffer and logs the request
- Sends back a hardcoded Redis-style response: `+OK\r\n`

## Installation

1. Make sure you have Go installed (version 1.16 or higher recommended)
2. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/redis-go.git
   cd redis-go
   ```
3. Build the project:
   ```bash
   go build
   ```

## Usage

1. Start the server:
   ```bash
   ./redis-go
   ```
   The server will start listening on `localhost:6379`

2. Connect to the server using a Redis client or `telnet`:
   ```bash
   telnet localhost 6379
   ```

## Project Structure

```
redis-go/
├── main.go          # Main server implementation
├── README.md        # Project documentation
└── go.mod          # Go module file
```

## Future Improvements

- [ ] Support multiple client connections
- [ ] Implement basic Redis commands (SET, GET, etc.)
- [ ] Add proper command parsing
- [ ] Implement data persistence
- [ ] Add configuration options
- [ ] Add unit tests
- [ ] Add proper error handling
- [ ] Add logging system

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.



