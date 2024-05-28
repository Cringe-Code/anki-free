# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY src/go.mod src/go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY ./src .

# Build the Go application
RUN go build -o anki

# Expose the application on port 8080
EXPOSE 8080

# Command to run the application
CMD ["./anki"]
