FROM golang:1.26-alpine AS builder
ARG VERSION=dev
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/api ./cmd/api

FROM alpine:3.20
RUN adduser -D -u 10001 appuser

COPY --from=builder /out/api /usr/local/bin/api
COPY configs /app/configs

WORKDIR /app
USER appuser
EXPOSE 8080

ENTRYPOINT ["api"]
