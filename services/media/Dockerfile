FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/media
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["media"]

EXPOSE 14000

# build
# docker build -t media .
# docker run --env-file=config.env --network media_backend -it --rm --name media-running media