# Use uma imagem base com Go
FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ordersystem ./cmd/ordersystem

FROM golang:1.22

WORKDIR /app

COPY --from=builder /app/ordersystem /app/ordersystem

COPY ./cmd/ordersystem/.env /app/.env
COPY ./app-run.sh /app/app-run.sh 

EXPOSE 8000 50051 8080

ENV ENV_FILE_PATH /app/.env

CMD ["/app/app-run.sh"]

