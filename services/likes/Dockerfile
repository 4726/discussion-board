FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/likes
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["likes"]

EXPOSE 14000

# docker build -t likes .
# docker run --env-file=config.env --network likes_backend -it --rm --name likes-running likes