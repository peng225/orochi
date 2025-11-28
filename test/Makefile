IMAGE_NAME := ghcr.io/peng225/orochi:latest
FLAGS := --build
N := 4
DATASTORES := $(addprefix data/datastore,$(shell seq 0 $(shell echo $$(($(N) - 1)))))

.PHONY: setup
setup:
	sudo apt-get install -y postgresql-client

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: test-io-only
test-io-only:
	go test -v -count=1 ./util.go ./object_test.go

data/datastore%:
	mkdir -p $@
	chmod 777 $@

.PHONY: prepare
prepare: $(DATASTORES)
	docker compose -p orochi up -d $(FLAGS)

.PHONY: clean
clean:
	docker compose -p orochi down
	rm -rf data
