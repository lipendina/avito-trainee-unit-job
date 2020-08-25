FROM golang:1.13.8

ADD . /avito
WORKDIR /avito

RUN go get github.com/google/uuid
RUN go get github.com/gorilla/mux
RUN go get github.com/jackc/pgx
RUN go get github.com/jackc/pgx/pgxpool
RUN go get golang.org/x/xerrors
RUN go get gopkg.in/yaml.v2

CMD ["go", "run", "main.go"]