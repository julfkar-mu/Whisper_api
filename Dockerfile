# Start from official Go image
FROM golang:1.21-alpine

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies layer
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy project files
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

# Build the server binary
RUN go build -o server ./cmd/server

# Expose server port
EXPOSE 8080

# Set environment variable default (can be overridden by docker run)
ENV OPENAI_API_KEY=""

# Command to run the server binary
CMD ["./server"]
