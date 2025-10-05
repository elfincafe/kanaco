FROM docker.io/library/golang:1.25.1-alpine3.22

ARG ROOT=/usr/local/kanaco

ENV TZ=Asia/Tokyo
ENV EDITOR=vi
ENV PS1="[\A]\u@kanaco:\W\\$ "
ENV STAGE=dev

RUN apk update && \
    apk upgrade && \
    apk add doas \
            vim \
            tzdata \
            gcc \
            make \
            musl-dev \
            cython \
            python3 \
            py3-pytest \
            cunit-dev && \
    rm -f /var/cache/apk/* && \
    adduser -u 2000 -D -h /home/dev -g users -G wheel dev && \
    echo -e "permit nopass :root\npermit nopass :wheel" > /etc/doas.conf && \
    mkdir -p ${ROOT}/src && \
    chown -R dev:users ${ROOT}

USER dev
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR ${ROOT}/src
