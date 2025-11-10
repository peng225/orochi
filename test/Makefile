IMAGE_NAME := ghcr.io/peng225/orochi:latest
FLAGS := --build

.PHONY: setup
setup:
	sudo apt-get install -y postgresql-client

.PHONY: test
test:
	go test -v ./...

data/datastore0:
	mkdir -p $@

data/datastore1:
	mkdir -p $@

data/datastore2:
	mkdir -p $@

.PHONY: prepare
prepare: data/datastore0 data/datastore1 data/datastore2
	docker compose -p orochi up -d $(FLAGS)
# FIXME: should wait until the manager container become ready.
	sleep 1
	PORTS="8082 8083 8084"; \
	for PORT in $${PORTS}; do \
		I=$$(($${PORT} - 8082)); \
		DATA="{\"baseURL\": \"http://datastore$${I}:$${PORT}\"}"; \
		curl -sS -X POST -H "Content-Type: application/json" \
			http://localhost:8080/datastores --data-raw "$${DATA}"; \
	done

.PHONY: clean
clean:
	docker compose -p orochi down
	rm -rf data
