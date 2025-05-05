FROM golang:1.22

WORKDIR /APP

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o curd-web-go .

EXPOSE 8000

CMD ["./curd-web-go"]


