FROM golang:1.16-alpine3.13
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 1234
CMD ["/app/main"]