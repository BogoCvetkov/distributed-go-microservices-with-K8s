# Build stuff
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go mod verify

COPY . .

WORKDIR /app/cmd/web

RUN go build -o ./build/fe-build

# Build tiny image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/web/build/fe-build /app
COPY --from=builder /app/cmd/web/templates /app/cmd/web/templates

# Don't run as root
RUN adduser -D -u 1001 bogo
RUN chown -R bogo:bogo /app
USER bogo

EXPOSE 80

CMD ["./fe-build"]