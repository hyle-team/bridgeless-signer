# Bridgeless signer

Bridgeless signer is a centralized service that performs verification of the
bridge deposit actions happening on the source chains and submit a signed
transaction to perform corresponding withdraw action on the target chain.

Although the service is built using GRPC and Protocol Buffers, it provides a
REST gateway to submit deposit transactions and check the status of the according withdraw.

## Build

Required environment:
```shell
export CGO_ENABLED=1
```

Build command (in repository root):
```shell
go build .
```

## Launch

### Configuration file 
Create a configuration file (`config.yaml`) with the following structure:

```yaml
log:
  level: debug
  disable_sentry: true

## PostgreSQL connection
db:
  url: postgres://signer:signer@signer-db/signer?sslmode=disable

## Port to listen for incoming GRPC requests
listener:
  addr: :8111

## Port to listen for incoming HTTP requests
rest_gateway:
  addr: :8222

## Available chains configuration
chains:
  list:
      ## Chain ID
    - id: "80002"
      ## RPC endpoint
      rpc:
      ## Bridge contract address
      bridge_address: "0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"
      ## Number of confirmations required for the deposit to be considered final
      confirmations: 1

## Service signer private key
signer:
  eth_signer: "8a8b8c...."
```

### Host environment:

Set up your host environment with the following variables:

```yaml
- name: KV_VIPER_FILE
  value: config.yaml # is the path to your config file
```

### Run service:
```shell
bridgeless-signer migrate up && bridgeless-signer run service
```

Or you can run the service with the following command:

```shell
make run
```

## Running with Docker Compose

Example of `docker-compose.yml` file:

```yml
services:
  signer-db:
    image: tokend/postgres-ubuntu:9.6
    hostname: signer-db
    container_name: signer-db
    restart: unless-stopped
    environment:
      - POSTGRES_USER=signer
      - POSTGRES_PASSWORD=signer
      - POSTGRES_DB=signer
      - PGDATA=/pgdata
    volumes:
      - signer-data:/pgdata
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U signer" ]
      interval: 5s
      timeout: 5s
      retries: 5
    ports:
      - 20000:5432

  signer:
    build:
      context: .
      dockerfile: Dockerfile
    hostname: signer
    container_name: signer
    restart: unless-stopped
    depends_on:
      signer-db:
        condition: service_healthy
    environment:
      KV_VIPER_FILE: '/config.yaml'
    volumes:
      - ./config.yaml:/config.yaml
    ports:
      - 8111:8111
      - 8222:8222
    entrypoint: sh -c "bridgeless-signer migrate up && bridgeless-signer run service"

volumes:
  signer-data:
```

## Tests

To run tests, execute the following command:

```shell
make test
```