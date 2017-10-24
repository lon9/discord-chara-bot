FROM golang:alpine

ADD . /go/src/github.com/Rompei/discord-chara-bot
WORKDIR /go/src/github.com/Rompei/discord-chara-bot


RUN apk add --no-cache git openssl
RUN wget https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 -O /usr/bin/dep
RUN chmod +x /usr/bin/dep
RUN dep ensure
RUN go install
