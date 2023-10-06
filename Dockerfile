# Stage 1 - Build the base
FROM golang:1.21.0-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project including nested Go files
COPY . .

# Build with CGo enabled
RUN GOOS=linux go build -o /psgc ./cmd/http

# Stage 2: Create a lightweight final image
FROM alpine:3.14.2

# Set the working directory
WORKDIR /app

# Set the environment variable
ENV ENV=prod

# Copy the binary built in the previous stage
COPY --from=builder /app/db ./db
COPY --from=builder /psgc .

# Run
CMD ["./psgc","api"]