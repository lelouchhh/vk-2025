# Dockerfile

FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o backend ./cmd/main.go

EXPOSE 8080

CMD ["./backend"]