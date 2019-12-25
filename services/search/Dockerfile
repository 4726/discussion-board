FROM golang:1.13.1

WORKDIR /go/src/github.com/4726/discussion-board/services/search
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o search .

CMD ["./search"]

EXPOSE 14000

# docker build -t search .
# docker run --env-file=config.env --network search_backend -it --rm --name search-running search