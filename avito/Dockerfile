FROM golang:1.13.8

ENV GO111MODULE=on

ADD . /avito
WORKDIR /avito

RUN go mod download

WORKDIR /avito/balance-service
RUN go build

WORKDIR /avito
CMD ["./balance-service/balance-service"]

#CMD ["go", "run", "main.go"]