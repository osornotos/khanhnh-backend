# Stage 1: Build the binary
FROM golang:1.23.8-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# RUN go test ./... -v

# Build the binary (build main.go from root directory)
RUN CGO_ENABLED=0 go build -o /app/backend .

# Stage 2: Final image
FROM alpine:3.20

RUN apk add --no-cache tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the builder stage
COPY --from=builder /app/backend /app/backend

# Copy additional files
COPY config.yaml /app/config.yaml
COPY migrations /app/migrations

# Set permissions for non-root user
RUN chown -R appuser:appgroup /app

WORKDIR /app
EXPOSE 8080

# Run as non-root user
USER appuser

ENTRYPOINT ["/app/backend"]