FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/posts/read
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["read"]

EXPOSE 14000

# docker build -t postsread . 
# docker run --env-file=config.env --network read_backend -it --rm --name postsread-running postsread