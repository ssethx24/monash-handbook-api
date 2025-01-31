FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o myapp

EXPOSE 8080

CMD ["./myapp"]
