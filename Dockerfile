FROM alpine:latest
MAINTAINER 1783296281@qq.com
RUN ["mkdir","/app"]
COPY src/main /app
EXPOSE 9999
CMD ["/app/main"]