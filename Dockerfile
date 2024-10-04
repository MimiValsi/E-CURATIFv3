FROM golang:1.21

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o main ./cmd

EXPOSE 3001

CMD ["./main"]
