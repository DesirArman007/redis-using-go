[![progress-banner](https://backend.codecrafters.io/progress/redis/b5d213e6-c7f2-4be8-9714-1b64ed7cd957)](https://app.codecrafters.io/users/DesirArman007?r=2qF)

# Build Your Own Redis (Go)

This repository contains my in-progress solution to the
[CodeCrafters Build Your Own Redis challenge](https://codecrafters.io/challenges/redis).
The challenge is to build a Redis server step by step, learning about TCP
servers, the Redis protocol (RESP), and the features Redis provides.

The implementation currently covers the early stages:

- `PING` - responds with `PONG`
- `ECHO <message>` - returns the message

This is not a complete Redis implementation. More commands and features will
be added as I progress through the challenge stages.

## Requirements

- Go 1.27

## Run locally

From the project directory, start the server:

```powershell
go run ./app
```

Or, in Bash or WSL, use the provided launcher:

```sh
./your_program.sh
```

The server listens on `127.0.0.1:6379`. Keep this terminal open while testing.
Press `Ctrl+C` to stop it.

## Test with redis-cli

In a second terminal, run:

```sh
redis-cli -h 127.0.0.1 -p 6379 PING
redis-cli -h 127.0.0.1 -p 6379 ECHO hello
```

Expected output:

```text
PONG
"hello"
```

You can also verify that the TCP port is accepting connections in PowerShell:

```powershell
Test-NetConnection 127.0.0.1 -Port 6379
```

## Run the test suite

```sh
go test ./...
```

