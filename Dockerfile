# Stage 1
FROM golang:1.25 AS builder

WORKDIR /work
COPY . .
RUN make

# Stage 2
FROM gcr.io/distroless/static

COPY --from=builder /work/orochi .
ENTRYPOINT [ "./orochi" ]
