FROM golang:1.23.2-alpine

WORKDIR /app

COPY . .

RUN go build -o main backend.go

ENV PORT = 8080

EXPOSE 8080

CMD ["./main"]