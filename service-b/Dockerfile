FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go build -o service main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/service .
CMD ["./service"]