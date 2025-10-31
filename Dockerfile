# Build stage
FROM golang:1.23-alpine AS builder

# Install git
RUN apk add --no-cache git ca-certificates

# Set build arguments
ARG GITHUB_USERNAME=mock-github-username
ARG GITHUB_TOKEN=mock-github-token


# Set environment variables for private module
RUN mkdir /user && \
    echo 'root:x:65534:65534:root:/:' > /user/passwd && \
    echo 'root:x:65534:' > /user/group
RUN printf "machine github.com login %s password %s" "$GITHUB_USERNAME" "$GITHUB_TOKEN" >> /root/.netrc
RUN chmod 600 /root/.netrc

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/app .
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"] 
