FROM golang:1.23.4

WORKDIR /app

COPY . .

RUN go build -o backend .

EXPOSE 8080

CMD ["./backend"]
