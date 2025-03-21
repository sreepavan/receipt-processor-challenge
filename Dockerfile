# Use the official Go image as the base image
FROM golang:1.24.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o receipt-processor .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/app/receipt-processor"]