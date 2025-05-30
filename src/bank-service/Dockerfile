FROM golang:1.24.2 AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN $GOPATH/bin/swag init -g ./cmd/server/main.go -o docs --parseDependency

RUN go test ./... -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bank-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/bank-service .
COPY config.env .
RUN chmod +x /root/bank-service

CMD ["./bank-service"]
