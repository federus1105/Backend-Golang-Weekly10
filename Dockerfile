FROM golang:1.25-alpine AS builder

WORKDIR /server

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /server

COPY --from=builder /server/main .

EXPOSE 8080

# Jalankan aplikasi.
CMD ["/server/main"]
