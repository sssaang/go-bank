# Build Stage
# build binary runnable file
FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
# copy the binary file to smaller env
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 1234
CMD ["/app/main"]