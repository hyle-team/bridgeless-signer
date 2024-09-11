# Bridgeless signer

Bridgeless signer is a centralized service that performs verification of the
bridge deposit actions happening on the source chains and submit a signed
transaction to perform corresponding withdraw action on the target chain.

Although the service is built using GRPC and Protocol Buffers, it provides a
REST gateway to submit deposit transactions and check the status of the according withdrawal.

## Architecture 

Service Core logic (`/internal/core`) consists of two parts: JSON API to submit deposits and check the 
status of the according withdrawal (`/core/api`), and different handlers to process withdrawals part by part,
where each part is being sent/consumed to/from RabbitMQ queues of messages (`/core/rabbitmq`).

There are two types of message consuming implemented here:
- Default consuming (`base`) - used to process incoming message immediately after consuming. Is used for
default scenarios where request can be processed independently;
- Batch consuming (`batch`) - used to collect incoming messages and process them after some period. Is used
for requests that should be better grouped and processed together to optimize/enhance processing (f. e. Bitcoin
transactions batching, Core tx submitting batching).

Multiple request handlers, that implement interface required by either base or
batch consumer, handle specific part of the withdrawal process (`/rabbitmq/consumer/processors`). They use
bridge processor (`/internal/bridge/processor`) to handle the request and then route the next one using 
RabbitMQ request producer (`/rabbitmq/producer`).

In order to avoid unexpected errors each request to process withdrawal can be resent up to
`maxCount` times in case of some system or third party services failure.

Bridge processor is the core system that interacts with database, Bridge Core module and chains for data 
retrieval/parsing/sending etc. It contains a set of proxies - implementations of the generalized method 
to work with different chains (`/bridge/proxy`). For now it supports EVM-based chains (`/proxy/evm`) and 
Bitcoin (`/proxy/btc`).

To better understand how withdrawals are processed by the service, lets look on the chain of requests processing:
1. Base/Batch Consumer reads pending request from the RabbitMQ queue;
2. Request is processed by specific consumer implementation;
3. Consumer implementation uses Bridge processor to handle request;
4. Bridge processor uses proxy to extract/parse/transform/send specific chain data;
5. After Bridge processor processed request, Consumer implementation forms and sends request to process the next
part of the withdrawal using Producer. Unexpectedly failed requests can also be resent by Base/Batch Consumer;
6. Returning to step #1 unless all parts of the withdrawal process are successfully finished.


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
      ## Chain type 
      type: evm
      ## RPC endpoint
      evm_rpc: "your_rpc_endpoint_here"
      bitcoin_rpc: 
        host: "your_rpc_endpoint_here" 
        user: "your_rpc_endpoint_here"
        pass: "your_rpc_endpoint_here"
      # bitcoin-specific data
      bitcoin_receivers: 
        - "list_of_addresses"
      # bitcoin-specific data
      network: testnet
      ## Bridge contract address
      bridge_address: "0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"
      ## Number of confirmations required for the deposit to be considered final
      confirmations: 1

## RabbitMQ configuration
rabbitmq:
  ## RabbitMQ connection URL
  url: amqp://guest:guest@localhost:5672/
  ## Number of instances per each consumer
  consumer_instances: 2
  ## Delivery resend parameters
  resend_params:
    ## Maximum number of retries
    max_retry_count: 5
    ## delivery resend delays
    delays: [1000, 2000, 5000, 10000, 20000, 60000]
  tx_submitter:
    max_size: 5
    period: 5s
  bitcoin_submitter:
    max_size: 5
    period: 10s


## Service signer private key
signer:
  eth_signer: "signer_private_key_here"
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

To pull the latest image of the service from the GitHub Container Registry, **firstly** execute the following command:

```shell
docker login ghcr.io
```

Example of `docker-compose.yml` file:

```yml
services:
  rabbitmq:
    image: rabbitmq:3.13.6-management-alpine
    hostname: rabbitmq
    container_name: 'rabbitmq'
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq
    ports:
      - 5672:5672
      - 15672:15672

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
    image: ghcr.io/hyle-team/bridgeless-signer:7108db395fe92c56875657190c4d9305376c4323
    hostname: signer
    container_name: signer
    restart: unless-stopped
    depends_on:
      signer-db:
        condition: service_healthy
      rabbitmq:
        condition: service_started
    environment:
      KV_VIPER_FILE: '/config.yaml'
    volumes:
      - ./config.docker.local.yaml:/config.yaml
    ports:
      - 8111:8111
      - 8222:8222
    entrypoint: sh -c "bridgeless-signer migrate up && bridgeless-signer run service"

volumes:
  rabbitmq-data:
  signer-data:
```

## Tests

To run tests, execute the following command:

```shell
make test
```