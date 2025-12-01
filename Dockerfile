FROM golang:1.23-alpine

WORKDIR /meet-book-api

COPY . .

# Install swag (pastikan git dan curl tersedia)
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    /go/bin/swag init -g cmd/server/main.go -o ./docs/

RUN go build -o meet-book-api cmd/server/main.go

EXPOSE 8080

CMD ["./meet-book-api"]

