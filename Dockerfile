# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY . .

# Download Go modules
RUN go mod download

# Build the Go application
RUN go build -o anki ./src && chmod +x ./anki

# Expose the application on port 8080
EXPOSE 8080

# Command to run the application
CMD ["./anki"]
