# Use the official Go image as the base image
FROM golang:1.21.4

# Set the working directory in the container
WORKDIR /app

# Copy the application files into the working directory
COPY . /app

# Build the application
RUN go build

# Expose port 8080
EXPOSE 8080

# Define the entry point for the container
CMD ["./main"]
