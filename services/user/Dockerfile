FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/user
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["user"]

EXPOSE 14000