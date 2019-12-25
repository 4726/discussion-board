FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/likes
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o likes .

CMD ["./likes"]

EXPOSE 14000

# docker build -t likes .
# docker run --env-file=config.env --network likes_backend -it --rm --name likes-running likes