IMAGE_NAME ?= ghcr.io/peng225/orochi

BINDIR := bin

GOLANGCI_LINT_VERSION := v2.5.0
GOLANGCI_LINT := $(BINDIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)


.PHONY: build
build:
	CGO_ENABLED=0 go build -o orochi -v main.go

.PHONY: image
image:
	docker build -t ghcr.io/peng225/orochi .

$(BINDIR):
	mkdir -p $@

.PHONY: generate
generate:
	go generate ./internal/...

$(GOLANGCI_LINT): | $(BINDIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b . $(GOLANGCI_LINT_VERSION)
	mv golangci-lint $(GOLANGCI_LINT)

.PHONY: lint
lint: | $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: test
test: build
	go test -v ./...

.PHONY: clean
clean:
	rm -f orochi
