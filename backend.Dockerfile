FROM golang:1.24-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build

#todo is it worth using go vendor? who does it work with private repositories?
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend .
RUN go build -ldflags="-s -w" -o /app/app cmd/main.go

FROM alpine AS runner

RUN addgroup -S appgroup \
    && adduser -S appuser -G appgroup

USER appuser

WORKDIR /home/appuser/app

COPY --from=builder /app/app app

ENTRYPOINT ["./app"]