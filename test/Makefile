IMAGE_NAME := ghcr.io/peng225/orochi:latest
FLAGS := --build
N := 4
DATASTORES := $(addprefix data/datastore,$(shell seq 0 $(shell echo $$(($(N) - 1)))))
FIRST_PORT := 8082
DATASTORE_PORTS := $(shell seq $(FIRST_PORT) $(shell echo $$(($(FIRST_PORT) + $(N) - 1))))
MANAGER_BASE_URL := http://localhost:8080

.PHONY: setup
setup:
	sudo apt-get install -y postgresql-client

.PHONY: test
test:
	go test -v ./...

data/datastore%:
	mkdir -p $@

.PHONY: prepare
prepare: $(DATASTORES)
	docker compose -p orochi up -d $(FLAGS)
# FIXME: should wait until the manager container become ready.
	sleep 1
	for PORT in $(DATASTORE_PORTS); do \
		I=$$(($${PORT} - $(FIRST_PORT))); \
		DATA="{\"baseURL\": \"http://datastore$${I}:$${PORT}\"}"; \
		curl -sS -X POST -H "Content-Type: application/json" \
			$(MANAGER_BASE_URL)/datastores --data-raw "$${DATA}"; \
	done

.PHONY: clean
clean:
	docker compose -p orochi down
	rm -rf data
