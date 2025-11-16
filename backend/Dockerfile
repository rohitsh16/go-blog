### Multi-stage Dockerfile for go-blog backend
FROM golang:1.23 as builder
WORKDIR /app

# Copy go.mod and download modules (cacheable)
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download

# Copy full repo
WORKDIR /app
COPY . .

# Build the binary
WORKDIR /app/backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-blog ./

## Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app/backend

# Copy binary and frontend assets
COPY --from=builder /app/go-blog /app/go-blog
COPY --from=builder /app/frontend /app/frontend
COPY --from=builder /app/backend/config_files/config.yaml /app/backend/config.yaml

EXPOSE 8080

CMD ["/app/go-blog"]
