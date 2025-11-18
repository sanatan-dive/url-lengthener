FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN ls -la && ls -la cmd/ && ls -la cmd/server/ || echo "cmd/server not found"
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-lengthener ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /url-lengthener /usr/local/bin/url-lengthener
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/url-lengthener"]
