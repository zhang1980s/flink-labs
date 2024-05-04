FROM golang:1.22 AS builder  

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o taxi-fare-generator

FROM alpine:latest  

WORKDIR /app
COPY --from=builder /code/taxi-fare-generator /app

ENTRYPOINT ["/app/taxi-fare-generator"]
