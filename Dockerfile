# Stage 1 - Build the base
FROM golang:1.21.0-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project including nested Go files
COPY . .

# Build
RUN apk add --no-cache gcc g++ #git openssh-client
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o psgc ./cmd/http

# Clean up unnecessary packages
RUN apk del gcc g++ #git openssh-client

# Stage 2: Create a lightweight final image
FROM alpine:3.14.2

# Set the working directory
WORKDIR /app

# Set the environment variable
ENV ENV=prod

# Copy the binary built in the previous stage
COPY --from=builder /app/db ./db
COPY --from=builder /app/psgc /usr/bin

# Run
CMD ["psgc", "api"]
