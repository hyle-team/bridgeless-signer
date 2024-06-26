.SILENT: gen-proto run

VIPER_FILE=config.yaml

gen-proto:
	cd proto && buf generate

run:
	KV_VIPER_FILE=$(VIPER_FILE) go run main.go run service
