FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/api /usr/local/bin/api
EXPOSE 8080
CMD ["/usr/local/bin/api"]
