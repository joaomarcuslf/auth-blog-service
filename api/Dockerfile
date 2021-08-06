FROM cosmtrek/air

WORKDIR /usr/app

COPY go.mod ./

RUN go mod download

RUN air -v

