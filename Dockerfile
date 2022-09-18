FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ./main

FROM alpine:latest
COPY --from=builder /app/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/main"]
