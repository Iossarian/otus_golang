FROM golang:1.18.3

WORKDIR /go/src

COPY . ${CODE_DIR}

RUN go test -i -tags integration ./tests/integration/...

CMD go test -v -tags integration ./tests/integration/...