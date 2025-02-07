# Dockerfile

FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o ping ./cmd/main.go

EXPOSE 8081

CMD ["./ping"]