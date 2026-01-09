## Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o git-stack-watch .

## Run stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY --from=builder /build/git-stack-watch /app/git-stack-watch

ENTRYPOINT ["/app/git-stack-watch"]
