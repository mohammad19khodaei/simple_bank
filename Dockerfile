# build stage
FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o main .

# run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 8080

CMD ["./main"]