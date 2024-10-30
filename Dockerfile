FROM golang:1.22

WORKDIR /app

COPY ./api .

RUN go mod download
RUN go build -o main .

CMD ["./main"]