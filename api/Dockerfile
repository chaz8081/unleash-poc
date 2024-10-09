FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./api .

RUN go build -o main .

CMD ["./main"]