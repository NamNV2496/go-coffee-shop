# Use the official Golang image as the base image
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go mod and sum files from the correct location
COPY ../../go.mod ../../go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod vendor

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/kitchen

######## Start a new stage from scratch #######
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
