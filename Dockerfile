FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM alpine:3.18 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/bin/api /usr/local/bin/api

COPY --from=builder /app/docs /app/docs

EXPOSE 8080
CMD ["/usr/local/bin/api"]
