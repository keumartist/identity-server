FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_TYPE=HTTP

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags $SERVICE_TYPE -o /identity-server ./cmd/server/main/http.go

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /identity-server /identity-server

EXPOSE 3000
CMD ["/identity-server"]
