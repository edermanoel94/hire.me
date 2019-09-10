FROM golang:latest

WORKDIR /go/src/desafio_bemobi

COPY . .

RUN go get -u github.com/golang/dep/...

RUN dep ensure

RUN go build -o api .

RUN go test -v ./...

CMD ["./api"]