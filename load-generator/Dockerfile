FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY load-generator ./load-generator
RUN go build -o /load-generator ./load-generator/cmd

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /load-generator /load-generator

USER nonroot:nonroot
ENTRYPOINT ["/load-generator"]
