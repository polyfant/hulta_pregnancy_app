# Use official Go image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project
COPY . .

# Copy React build
COPY frontend-react/build ./frontend/build

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy built binary and static files
COPY --from=builder /app/main .
COPY --from=builder /app/frontend/build ./frontend/build

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
