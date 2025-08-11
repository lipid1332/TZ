FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app .

EXPOSE ${APP_PORT}

CMD ["./myapp"]
