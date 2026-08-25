FROM golang:1.26-alpine AS builder

ENV GOCACHE=/root/.cache/go-build

RUN apk add --no-cache build-base

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN sqlc generate

RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 go build \
        -ldflags="-s -w" \
        -o govd ./cmd/main.go

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache ffmpeg

COPY --from=builder /app/govd ./govd

ENTRYPOINT ["./govd"]