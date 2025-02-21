# Build stage
FROM golang:1.23.6 AS builder
WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

# Copy the entire project
COPY . .

# Move into the api/ directory
RUN cd api && go mod download

# Final stage
FROM golang:1.23.6
WORKDIR /app

# Copy the project files
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY . .

# Set the working directory to api/
WORKDIR /app/api

# Run Air for hot reloading
CMD ["air"]