.SILENT: gen-proto run test

TESTING_PACKAGES = ./internal/bridge/processor/...
VIPER_FILE="${GOPATH}/src/github.com/hyle-team/bridgeless-signer/config.local.yaml"


gen-proto:
	cd proto && buf generate

build:
	rm -f $(GOPATH)/bin/signer
	go build -o $(GOPATH)/bin/signer

run:
	build
	KV_VIPER_FILE=$(VIPER_FILE) signer migrate up
	KV_VIPER_FILE=$(VIPER_FILE) signer run service

test:
	KV_VIPER_FILE=$(VIPER_FILE) go test -count=1 $(TESTING_PACKAGES)
