FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/media
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o media .

CMD ["./media"]

EXPOSE 14000

# build
# docker build -t media .
# docker run --env-file=config.env --network media_backend -it --rm --name media-running media


