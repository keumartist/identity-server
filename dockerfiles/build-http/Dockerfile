FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_TYPE=HTTP
# 정적 바이너리 생성을 위해 CGO_ENABLED=0 설정
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags $SERVICE_TYPE -o /identity-server ./cmd/server/main/http.go

FROM scratch
COPY --from=builder /identity-server /identity-server

EXPOSE 8080
CMD [ "/identity-server" ]
