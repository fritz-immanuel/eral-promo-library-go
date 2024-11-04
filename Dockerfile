# Use the official Go image to build the application
FROM golang:1.21 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o eral-promo-library-go

# Use a minimal base image to reduce the size of the container
FROM alpine:latest

# Set the PORT environment variable to 8080 for Cloud Run
ENV PORT=8080

# Copy the built application binary from the builder image
COPY --from=builder /app/eral-promo-library-go /eral-promo-library-go

# Expose port 8080 for Cloud Run (internal routing port)
EXPOSE 8080

# Start the application and use the PORT environment variable
CMD ["sh", "-c", "/eral-promo-library-go --port=$PORT"]
