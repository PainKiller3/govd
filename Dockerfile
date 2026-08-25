FROM golang:1.24-alpine3.21 AS builder

ENV GOCACHE=/root/.cache/go-build
ENV GOTOOLCHAIN=auto

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    --mount=type=cache,target=/var/lib/apk,sharing=locked \
    apk add --no-cache \
        --repository="https://dl-cdn.alpinelinux.org/alpine/v3.21/main" \
        --repository="https://dl-cdn.alpinelinux.org/alpine/v3.21/community" \
        build-base \
        libheif-dev

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN sqlc generate

RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=1 go build \
        -ldflags="-s -w" \
        -o govd ./cmd/main.go

FROM alpine:3.21 AS runtime

WORKDIR /app

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    --mount=type=cache,target=/var/lib/apk,sharing=locked \
    apk add --no-cache \
        --repository="https://dl-cdn.alpinelinux.org/alpine/v3.21/main" \
        --repository="https://dl-cdn.alpinelinux.org/alpine/v3.21/community" \
        ffmpeg \
        libheif

COPY --from=builder /app/govd ./govd

ENTRYPOINT ["./govd"]