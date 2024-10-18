# Use the official Golang image as a build stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o model-viewer .

# Final stage: minimal base image to run the application
FROM alpine:3.20

# Install certificates to make HTTPS requests work if needed
RUN apk --no-cache add ca-certificates openssl

# Set the working directory
WORKDIR /app

ENV DEVELOPMENT=false

# Copy the built application from the builder stage
COPY --from=builder /app/model-viewer /app/

# Copy the embedded templates
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

RUN openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=localhost"

# Expose the port that the application will listen on
EXPOSE 80
EXPOSE 443

# Command to run the application
CMD ["/app/model-viewer"]
