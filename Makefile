.SILENT: gen-proto run test

TESTING_PACKAGES = ./internal/bridge/evm...
VIPER_FILE=config.yaml

# set the rpc url for the sepolia network
export SEPOLIA_RPC_URL=https://warmhearted-green-...


gen-proto:
	cd proto && buf generate

run:
	KV_VIPER_FILE=$(VIPER_FILE) go run main.go run service

test:
	go test -count=1 $(TESTING_PACKAGES)
