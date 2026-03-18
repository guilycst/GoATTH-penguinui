FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /src /app

EXPOSE 8090

CMD ["/app/server", "-port", "8090"]
