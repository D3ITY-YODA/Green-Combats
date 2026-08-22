FROM golang:1.26-alpine AS builder
ARG VERSION=dev
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/api ./cmd/api

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w" \
    -o /out/migrate ./cmd/migrate

FROM alpine:3.20
RUN adduser -D -u 10001 appuser

COPY --from=builder /out/api /usr/local/bin/api
COPY --from=builder /out/migrate /usr/local/bin/migrate
COPY configs /app/configs
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

WORKDIR /app
USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
