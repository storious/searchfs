# ZigKV

ZigKV is a learning-oriented in-memory key-value cache written in Zig.

It is part of SystemLab and focuses on low-level systems programming topics such as memory management, command parsing, cache behavior, and storage internals.

ZigKV is not a production Redis replacement. It is a systems programming lab.

## Current Features

- In-memory key-value store
- `SET`, `GET`, `DEL`, `EXISTS`, `SETEX`, `PING`
- TTL support with logical clock input
- Command parser
- Engine layer
- Response formatting
- Unit tests
- Makefile integration
- CI integration

## Architecture

```text
command.zig   parse text commands
store.zig     manage key-value data and TTL
engine.zig    execute commands against the store
response.zig  format protocol responses
clock.zig     provide time abstraction
main.zig      executable entrypoint
```

## Commands

| Command | Description |
|---------|-------------|
| `PING` | Check whether the engine is alive |
| `SET key value` | Set a key |
| `GET key` | Get a key |
| `DEL key` | Delete a key |
| `EXISTS key` | Check whether a key exists |
| `SETEX key ttl_ms value` | Set a key with TTL in milliseconds |
| `TTL key` | Get remaining TTL |
| `PERSIST key` | Remove TTL from a key |
| `KEYS` | List all live keys |
| `DBSIZE` | Return number of live keys |
| `CLEAR` | Remove all keys |

## Current Features

- In-memory key-value store
- Logical TTL and expiration cleanup
- Command parser
- Engine execution layer
- Protocol response formatting
- Core/runtime separation
- Single-shot CLI runtime
- Unit tests
- Makefile and CI integration
