FROM golang:alpine3.15

WORKDIR /app

COPY . .

RUN go get

RUN go build -o main .

CMD ["go", "run", "main.go"]