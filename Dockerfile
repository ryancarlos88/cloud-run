# Stage 1: Build the Go binary
FROM golang:1.20 AS builder

WORKDIR /app

# Copy the Go module files
COPY go.mod ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal runtime image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/app .

EXPOSE 8080
# Set the entrypoint for the container
ENTRYPOINT ["./app"]