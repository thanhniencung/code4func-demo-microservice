FROM golang:1.12-alpine

RUN apk update && apk add git

ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/user-service

COPY . .

RUN go mod init user-service
WORKDIR $GOPATH/src/user-service

RUN GOOS=linux go build -o user-service

ENTRYPOINT ["./user-service"]

EXPOSE 3000