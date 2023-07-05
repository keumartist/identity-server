# Go server application to serve identity service

This project is an identity service implemented in Go. It's designed to provide user authentication and authorization for other related applications. The service supports HTTP interfaces for general auth and gRPC interfaces for token verification

---

## Tech stack

Language: Go (1.20)

Framework: Fiber (for HTTP server), gRPC (for RPC server)

Testing: Testify

---

## How to run

Set environment variables

```shell
export PRIVATE_KEY_PATH=<your_private_key_path>
export SECRET_KEY_PATH=<your_secret_key_path>
export ISSUER=<your_issuer>
export MYSQL_CONNECTION_STRING=<your_mysql_connection_string>
export GOOGLE_OAUTH_CLIENT_ID=<your_google_oauth_client_id>
export GOOGLE_OAUTH_CLIENT_SECRET=<your_google_oauth_client_secret>
export AWS_REGION=<your_aws_region>
export SENDER_EMAIL=<your_sender_email>
```

Install dependencies

```shell
go mod download
```

Run GRPC server

```shell
go run cmd/server/main/grpc.go
```

Run HTTP server

```shell
go run cmd/server/main/http.go
```

---

## How to run the tests

```shell
go test ./...
```

---

## How to build using docker

Build image for GRPC app

```shell
docker build -f dockerfiles/build-grpc/Dockerfile .
```

Build image for HTTP app

```shell
docker build -f dockerfiles/build-http/Dockerfile .
```
