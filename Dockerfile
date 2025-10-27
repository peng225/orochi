# syntax=docker/dockerfile:1.4

# Stage 1
FROM golang:1.25 AS builder

WORKDIR /work

# Copy dependencies first for cache.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o orochi -v main.go

# Stage 2
FROM gcr.io/distroless/static

COPY --from=builder /work/orochi .
ENTRYPOINT ["./orochi"]
