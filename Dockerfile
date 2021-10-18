FROM golang:1.17.1
MAINTAINER 1783296281@qq.com
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=80
WORKDIR /go/src/
COPY src/main.go .
RUN go build main.go
EXPOSE 9999
ENTRYPOINT ["./main"]