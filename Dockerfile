FROM golang:latest AS builder
WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bsky-bot-app ./

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /builder/bsky-bot-app ./bsky-bot-app
EXPOSE 8000
CMD ["./bsky-bot-app"]
