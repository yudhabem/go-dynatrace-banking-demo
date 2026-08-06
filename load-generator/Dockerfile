FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o load-generator ./cmd

CMD ["./load-generator"]
