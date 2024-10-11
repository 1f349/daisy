.PHONY: build dev

all:
	@echo Nothing
	cd web && yarn

build:
	cd web && yarn build
	mkdir -p dist && go build -o dist/daisy ./cmd/daisy

dev:
	cd web && yarn build
	go run -tags nullauth ./cmd/daisy serve -conf .data/config.json
