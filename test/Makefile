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
	docker network create testnet
	docker run -d \
		--name postgres \
		--network testnet \
		-e POSTGRES_USER=testuser \
		-e POSTGRES_PASSWORD=testpass \
		-e POSTGRES_DB=orochi \
		-p 8180:5432 \
  		postgres
	until pg_isready -h localhost -p 8180 -U postgres > /dev/null 2>&1; do \
  		echo "Waiting for PostgreSQL to be ready..."; \
  		sleep 1; \
	done
	psql "postgresql://testuser:testpass@localhost:8180/orochi" -f ../internal/manager/infra/postgresql/sqlc/schema.sql
	docker run -d \
		--name manager \
		--network testnet \
		-e DSN="host=postgres port=5432 user=testuser password=testpass dbname=orochi sslmode=disable" \
		-p 8080:8080 \
		$(IMAGE_NAME) \
		manager
	docker run -d --name datastore0 --network testnet -p 8082:8082 $(IMAGE_NAME) datastore
	docker run -d --name datastore1 --network testnet -p 8083:8083 $(IMAGE_NAME) datastore
	docker run -d --name datastore2 --network testnet -p 8084:8084 $(IMAGE_NAME) datastore
	docker run -d \
		--name gateway \
		--network testnet \
		-e DSN="host=postgres port=5432 user=testuser password=testpass dbname=orochi sslmode=disable" \
		-p 8081:8081 $(IMAGE_NAME) \
		gateway

.PHONY: clean
clean:
	docker stop gateway || true
	docker rm gateway || true
	docker stop datastore0 || true
	docker rm datastore0 || true
	docker stop datastore1 || true
	docker rm datastore1 || true
	docker stop datastore2 || true
	docker rm datastore2 || true
	docker stop manager || true
	docker rm manager || true
	docker stop postgres || true
	docker rm postgres || true
	docker network rm testnet || true
