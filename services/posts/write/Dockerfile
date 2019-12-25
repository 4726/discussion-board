FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/posts/write
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o write .

CMD ["./write"]

EXPOSE 14000

# docker build -t postswrite . 
# docker run --env-file=config.env --network write_backend -it --rm --name postswrite-running postswrite