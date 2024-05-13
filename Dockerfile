# Use the official Golang image from Docker Hub
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Copy go.mod and go.sum files to the working directory
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/journie

# # Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD ["./main"]
