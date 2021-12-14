FROM alpine:latest
MAINTAINER "tscuite@qq.com"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update --no-cache && apk add ca-certificates
WORKDIR /app
#COPY eventrouter ./
#RUN ls
#ENTRYPOINT ./eventrouter -v 3 -logtostderr
COPY main ./
ENTRYPOINT ./main
