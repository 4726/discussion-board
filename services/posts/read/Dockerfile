FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/posts/read
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o read .

CMD ["./read"]

EXPOSE 14000

# docker build -t postsread . 
# docker run --env-file=config.env --network read_backend -it --rm --name postsread-running postsread