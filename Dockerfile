FROM golang as build

ENV GOPROXY=https://goproxy.io

ADD . /chatroom

WORKDIR /chatroom/cmd/chatroom

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatroom

FROM alpine:3.7

RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates && \
    echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /www/conf

WORKDIR /www

COPY --from=build /chatroom/cmd/chatroom/chatroom /usr/bin/chatroom

RUN chmod +x /usr/bin/chatroom

ENTRYPOINT ["chatroom"]
