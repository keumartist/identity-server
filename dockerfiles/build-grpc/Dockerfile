FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_TYPE=GRPC
RUN go build -tags $SERVICE_TYPE -o /identity-server ./cmd/server/main/grpc.go

FROM scratch
COPY --from=builder /identity-server /identity-server

EXPOSE 50051
CMD [ "/identity-server" ]
