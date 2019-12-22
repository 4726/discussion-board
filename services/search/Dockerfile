FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/search
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["search"]

EXPOSE 14000

# docker build -t search .
# docker run --env-file=config.env --network search_backend -it --rm --name search-running search