FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_TYPE=HTTP
RUN go build -tags $SERVICE_TYPE -o /identity-server ./cmd/server/main/http.go

FROM scratch
COPY --from=builder /identity-server /identity-server

EXPOSE 8080
CMD [ "/identity-server" ]
