FROM docker.io/library/alpine:3.21.2

ENV TZ Asia/Tokyo
ENV EDITOR vi
ENV GOPATH /usr/local/kanaco
ENV PATH ${GOPATH}/bin:${PATH}
ENV PS1 "[\A]\u@kanaco:\W\\$ "
ENV STAGE dev

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
    echo -e "permit nopass :root\npermit nopass :wheel" > /etc/doas/d/doas.conf && \
    mkdir -p ${GOPATH}/bin ${GOPATH}/pkg && \
    chown -R dev:users ${GOPATH}/bin ${GOPATH}/pkg

USER dev
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR ${GOPATH}/src
