FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /phoenix-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata python3 py3-pip
COPY --from=builder /phoenix-api /usr/local/bin/phoenix-api
COPY configs/ /etc/phoenix/configs/
WORKDIR /app
EXPOSE 8066
CMD ["phoenix-api"]
