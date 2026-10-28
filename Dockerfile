ARG GO_VERSION=1.27.3
FROM golang:${GO_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/migrate ./cmd/migration

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S goapp \
    && adduser -S -G goapp goapp

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations

USER goapp

EXPOSE 8080

CMD ["/app/api"]
