FROM alpine:3.16

LABEL version="0.1.0"
LABEL description="KanaCo Kana Converter Library"
LABEL maintainer="Takamichi Seki <takamichi.seki@elfincafe.com>"

RUN apk update && apk upgrade && apk add doas vim tzdata gcc make musl-dev go go-doc cython cython-doc python3 python3-doc
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
RUN apk del tzdata
RUN rm -f /var/cache/apk/*

RUN addgroup -g 50 staff
RUN adduser -u 2000 -D -h /home/zlkv -G staff zlkv
RUN addgroup zlkv wheel
RUN echo "permit nopass :root" >> /etc/doas.d/doas.conf
RUN echo "permit nopass :wheel" >> /etc/doas.d/doas.conf
ENV GOPATH /var/opt
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/pkg" "$GOPATH/bin" && chown -R zlkv:staff "$GOPATH/pkg" "$GOPATH/bin"
WORKDIR $GOPATH/src
RUN touch ~/.profile
RUN echo 'export PS1="[\A]\u@kanaco:\W\\$ "' >> ~/.profile
USER zlkv
RUN touch ~/.profile
RUN echo 'export PS1="[\A]\u@kanaco:\W\\$ "' >> ~/.profile
