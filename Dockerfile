# Step 1: Start from a lightweight Go image
FROM golang:1.23-alpine AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files for dependency installation
COPY go.mod go.sum ./

# Step 4: Download Go dependencies
RUN go mod download

# Step 5: Copy the rest of the application source code
COPY . .

# Step 6: Build the Go application
RUN go build -o main .

# Step 7: Use a smaller image to run the Go application
FROM alpine:latest

# Step 8: Set the working directory for the final container
WORKDIR /root/

# Step 9: Copy the binary from the builder container to the new one
COPY --from=builder /app/main .

# Step 10: Expose the port that your Go application uses
EXPOSE 8080

# Step 11: Set the entry point to run your Go application
CMD ["./main"]
