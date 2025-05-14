# RedGo

RedGo is a lightweight in-memory key-value store written in Go, inspired by Redis. It provides basic Redis-like functionality and supports interaction via the official `redis-cli` or a custom client (`redgo-cli`). The project is designed to be simple, efficient, and easy to use.

## Features

- **Redis Compatibility**: Fully compatible with `redis-cli` for command execution.
- **Key-Value Store**: Supports basic commands like `SET`, `GET`, and `DEL`.
- **Hash Operations**: Commands like `HSET`, `HGET`, and `HGETALL` for managing hashes.
- **Pub/Sub**: Publish/subscribe functionality for real-time messaging.
- **Persistence**: Append-only file (AOF) support for data durability.
- **Custom Client**: Includes a custom CLI (`redgo-cli`) for interacting with the server.

## Project Structure

```
redgo/
├── LICENSE                # License file
├── redgo-cli/             # Custom client implementation
│   ├── cli.go             # CLI logic
│   ├── go.mod             # Module dependencies
│   ├── go.sum             # Dependency checksums
│   ├── main.go            # Entry point for the client
│   ├── parser.go          # Command parsing logic
│   └── value_types.go     # Data type definitions
├── redgo-server/          # Server implementation
│   ├── aof.go             # Append-only file (AOF) persistence
│   ├── database.aof       # AOF data file
│   ├── go.mod             # Module dependencies
│   ├── handler.go         # Command handler logic
│   ├── hash.go            # Hash command implementations
│   ├── main.go            # Entry point for the server
│   ├── parser.go          # Command parsing logic
│   ├── pub_sub.go         # Pub/Sub functionality
│   ├── string.go          # String command implementations
│   └── value_types.go     # Data type definitions
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/theZeuses/redgo.git
   cd redgo
   ```

2. Build the server:
   ```bash
   cd redgo-server
   go build -o redgo-server
   ```

3. Build the custom client (optional):
   ```bash
   cd ../redgo-cli
   go build -o redgo-cli
   ```

## Running the Server

1. Start the server:
   ```bash
   ./redgo-server
   ```

2. The server will start listening on port `7000` by default.

## Using the Server

### With `redis-cli`

You can use the official `redis-cli` to interact with the server:

1. Install `redis-cli` (if not already installed):
   ```bash
   sudo apt install redis-tools
   ```

2. Connect to the server:
   ```bash
   redis-cli -p 7000
   ```

3. Use commands like:
   - `SET key value`
   - `GET key`
   - `DEL key`
   - `HSET myhash field1 value1`
   - `HGET myhash field1`
   - `HGETALL myhash`
   - `PUBLISH channel message`
   - `SUBSCRIBE channel`

### With the Custom Client

Alternatively, you can use the custom client:

1. Run the client:
   ```bash
   ./redgo-cli
   ```

2. Enter commands interactively, similar to `redis-cli`.

## Available Commands

### String Commands
- **SET key value**: Set a key to a value.
- **GET key**: Get the value of a key.

### Hash Commands
- **HSET key field value**: Set a field in a hash.
- **HGET key field**: Get the value of a field in a hash.
- **HGETALL key**: Get all fields and values in a hash.

### Pub/Sub Commands
- **PUBLISH channel message**: Publish a message to a channel.
- **SUBSCRIBE channel**: Subscribe to a channel to receive messages.

## Persistence

The server supports append-only file (AOF) persistence. All write operations are logged to the `database.aof` file, ensuring data durability across restarts.

The cli retain command history across sessions.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the project.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.