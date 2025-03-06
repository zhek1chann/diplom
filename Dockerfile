FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /bin/crud_server .

CMD ["./crud_server"]
