FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /url-shortener

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/url-shortener ./cmd/main.go

FROM alpine:3.23.3

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /usr/local/bin/url-shortener .

EXPOSE 8080

ARG STORAGE_TYPE
ENV STORAGE_TYPE "$STORAGE_TYPE"

CMD ["./url-shortener"]
