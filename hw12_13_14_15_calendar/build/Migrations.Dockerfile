FROM golang:1.16.2

WORKDIR "/migrations"

COPY ./migrations .

RUN go get -u github.com/pressly/goose/cmd/goose

CMD ["/go/bin/goose", "postgres", "postgres://calendar:calendar@calendar_db:5432/calendar?sslmode=disable", "up"]
