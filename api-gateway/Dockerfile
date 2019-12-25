FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/api-gateway
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gateway .

CMD ["./gateway"]

EXPOSE 14000

# build
# docker build -t api-gateway .
# docker run --env-file=config.env -it --rm --name gateway-running api-gateway 