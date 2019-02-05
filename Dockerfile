FROM golang:1.11-alpine

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /app

ENV GO111MODULE=on

COPY ./src/go.mod .
COPY ./src/go.sum .

RUN go mod download

COPY ./src /app

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/app 

CMD ["/bin/app"]