VERSION 0.6
FROM golang:1.19
WORKDIR /go-workdir

build-http-echo-uds:
    COPY ./http_server_echo_uds/main.go .
    RUN go build -o /app/my-app main.go
    SAVE ARTIFACT /app/my-app

build-simple-http:
    COPY ./simple_http/main.go .
    RUN go build -o /app/my-app main.go
    SAVE ARTIFACT /app/my-app

docker-simple-http-uds:
    COPY +build-http-echo-uds/my-app .
    ENTRYPOINT ["/go-workdir/my-app"]
    SAVE IMAGE simple-http-uds:latest

docker-simple-http:
    COPY +build-simple-http/my-app .
    ENTRYPOINT ["/go-workdir/my-app"]
    SAVE IMAGE simple-http:latest

build-all:
    BUILD +docker-simple-http
    BUILD +docker-simple-http-uds
