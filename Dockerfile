ARG GO_VERSION=1.24.5

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add --no-cache build-base git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -o app ./cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
