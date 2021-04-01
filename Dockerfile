FROM alpine:3.13

LABEL version="0.1.0"
LABEL description="KanaCo Laibrary for Go"

ENV GO111MODULE=on
ENV CC=clang
ENV EDITOR=vi
ENV PAGER=less
ENV PS1='[\t]\u@\h:\W\\$ '

RUN apk update
RUN apk upgrade
RUN apk add go vim tzdata
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
RUN apk del tzdata
RUN rm -f /var/cache/apk/*

