FROM docker.io/library/alpine:3.21.2

ENV GOPATH /usr/local/kanaco
ENV PATH ${GOPATH}/bin:${PATH}
ENV PS1 "[\A]\u@kanaco:\W\\$ "

RUN apk update && \
    apk upgrade && \
    apk add doas \
            vim \
            tzdata \
            clang \
            make \
            musl-dev \
            go \
            cython \
            python3 \
            py3-pytest \
            cunit-dev && \
    rm -f /var/cache/apk/* && \
    adduser -u 2000 -D -h /home/dev -G users dev && \
    addgroup dev wheel && \
    echo -e "permit nopass :root\npermit nopass :wheel" > /etc/doas.conf && \
    mkdir -p ${GOPATH}/pkg ${GOPATH}/bin ${GOPATH}/src && \
    chown -R dev:users ${GOPATH}

USER dev
WORKDIR ${GOPATH}/src
