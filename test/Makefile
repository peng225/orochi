IMAGE_NAME := ghcr.io/peng225/orochi:latest

.PHONY: setup
setup:
	sudo apt-get install -y postgresql-client

.PHONY: test
test: clean
	$(MAKE) prepare
	go test -v ./...
	$(MAKE) clean

.PHONY: prepare
prepare:
	docker compose -p orochi up -d --build

.PHONY: clean
clean:
	docker compose -p orochi down
