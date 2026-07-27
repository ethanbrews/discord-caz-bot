FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static, lightweight binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o discord_caz_bot .

# Final lean stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/discord_caz_bot .

CMD ["./discord_caz_bot"]