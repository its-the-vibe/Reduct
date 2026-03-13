# Reduct

A Go service that fans-in multiple Redis Pub/Sub source channels into a single
target channel.

## How it works

Reduct subscribes to one or more Redis source channels. Every message received
on any of those channels is re-published verbatim to a configured target
channel. This lets you aggregate streams from many producers into one consumer.

```
source-channel-1 ──┐
source-channel-2 ──┤──▶ reduct ──▶ output-channel
source-channel-N ──┘
```

## Prerequisites

- [Go 1.24+](https://go.dev/dl/) (for running locally)
- [Docker](https://docs.docker.com/get-docker/) and
  [Docker Compose](https://docs.docker.com/compose/) (for containerised
  deployment)
- An externally hosted Redis server (≥ 6.0 recommended)

## Configuration

### 1. Config file

Copy the example config and edit it to match your setup:

```bash
cp config.example.yaml config.yaml
```

`config.yaml` (gitignored — never committed):

```yaml
redis:
  addr: "your-redis-host:6379"
  db: 0

channels:
  target: "output-channel"
  sources:
    - "source-channel-1"
    - "source-channel-2"
```

### 2. Environment variables (`.env`)

Sensitive values are read from environment variables.
Copy the example file and set the real password:

```bash
cp .env.example .env
# edit .env and replace the placeholder value
```

`.env` (gitignored — never committed):

```
REDIS_PASSWORD=your-redis-password-here
```

When running inside Docker Compose, the `env_file` directive automatically
injects these variables into the container.

## Running locally

```bash
go run . # reads config.yaml and .env from the current directory
```

Override the config path with the `CONFIG_FILE` environment variable:

```bash
CONFIG_FILE=/path/to/config.yaml go run .
```

## Running with Docker Compose

```bash
# Build the image
docker compose build

# Start the service
docker compose up -d

# View logs
docker compose logs -f reduct

# Stop the service
docker compose down
```

The container runs as **read-only** (`read_only: true`).
`config.yaml` is bind-mounted into the container at `/config.yaml`.
`.env` is loaded automatically by Docker Compose.

## Building the Docker image manually

```bash
docker build -t reduct:latest .
docker run --rm \
  --env-file .env \
  -e CONFIG_FILE=/config.yaml \
  -v "$(pwd)/config.yaml:/config.yaml:ro" \
  --read-only \
  reduct:latest
```
