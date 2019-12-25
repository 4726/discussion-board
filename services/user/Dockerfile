FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/user
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user .

CMD ["./user"]

EXPOSE 14000