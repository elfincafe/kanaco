FROM docker.io/library/alpine:3.19.1

LABEL version="0.1.0"
LABEL description="KanaCo Kana Converter Library"
LABEL maintainer="Takamichi Seki <takamichi.seki@elfincafe.com>"

ENV GOPATH /usr/local/kanaco
ENV PATH ${GOPATH}/bin:${PATH}
ENV PS1 "[\A]\u@kanaco:\W\\$ "

RUN apk update && \
    apk upgrade && \
    apk add doas \
            vim \
            tzdata \
            gcc \
            make \
            musl-dev \
            go \
            cython \
            python3 \
            py3-pytest \
            cunit-dev && \
    rm -f /var/cache/apk/* && \
    adduser -u 2000 -D -h /home/zlkv -G users zlkv && \
    addgroup zlkv wheel && \
    echo -e "permit nopass :root\npermit nopass :wheel" > /etc/doas.d/doas.conf && \
    mkdir -p ${GOPATH}/pkg ${GOPATH}/bin ${GOPATH}/src && \
    chown -R zlkv:users ${GOPATH}

USER zlkv
WORKDIR $GOPATH/src
