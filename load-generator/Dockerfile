FROM golang:1.25

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o load-generator ./cmd

CMD ["./load-generator"]