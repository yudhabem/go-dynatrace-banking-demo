FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/gateway

RUN go build -o gateway ./cmd

EXPOSE 8080

CMD ["./gateway"]