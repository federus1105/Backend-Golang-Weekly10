# Stage 1: Build Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod/sum and download deps first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN go build -o server ./cmd/main.go

# Stage 2: Run the binary from a minimal image
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]
