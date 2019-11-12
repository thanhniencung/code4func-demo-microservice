FROM golang:1.12-alpine

RUN apk update && apk add git

ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/order-service

COPY . .

RUN go mod init order-service
WORKDIR $GOPATH/src/order-service

RUN GOOS=linux go build -o order-service

ENTRYPOINT ["./order-service"]

EXPOSE 3001