FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/app cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/app /app/app

ENTRYPOINT ["./app"]