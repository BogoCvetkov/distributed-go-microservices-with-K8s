# Build stuff
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go mod verify

COPY . .

WORKDIR /app/cmd/api

RUN go build -o ./build/logger-build

# Build tiny image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/api/build/logger-build /app

# Don't run as root
RUN adduser -D -u 1001 bogo
RUN chown -R bogo:bogo /app
USER bogo

EXPOSE 3002

CMD ["./logger-build"]