# Docker Usage Guide for GoInspectorGadget

This guide explains how to use Docker with the GoInspectorGadget application.

## Prerequisites

- Docker installed on your system
- Docker Compose installed on your system

## Basic Usage

### Building the Docker Image

```bash
docker build -t goinspectorgadget .
```

### Running the Container (One-Off Command)

```bash
docker run --rm goinspectorgadget [command]
```

Examples:
```bash
# Show help
docker run --rm goinspectorgadget help

# Create a new case
docker run --rm -v $(pwd)/data:/app/data goinspectorgadget case create --title "Murder Case" --desc "Investigation of murder on Main St" --type "Homicide"
```

## Using Docker Compose

### Standard Usage (Run and Exit)

```bash
docker-compose up
```

This will run the container with the default "help" command and then exit.

### Development Mode (Keep Container Running)

```bash
# Start the container in the background
docker-compose -f docker-compose.dev.yml up -d

# Execute commands against the running container
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app [command]

# Examples:
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app case list
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app evidence list

# Stop the container when done
docker-compose -f docker-compose.dev.yml down
```

## Data Persistence

The application stores data in the `/app/data` directory inside the container. This is mapped to the `./data` directory in your project folder.

## Common Commands

```bash
# Show help
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app help

# Create a new case
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app case create --title "New Case" --desc "Description" --type "Fraud"

# List all cases
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app case list

# Add evidence to a case
docker exec -it goinspectorgadget-goinspectorgadget-1 /app/app evidence add --desc "Fingerprint" --type "PHYSICAL" --case [case-id]
```

## Troubleshooting

If you encounter issues:

1. Check container logs: `docker logs goinspectorgadget-goinspectorgadget-1`
2. Verify volume mounts: `docker inspect goinspectorgadget-goinspectorgadget-1`
3. Ensure data directory has correct permissions: `chmod 777 ./data` 