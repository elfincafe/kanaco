FROM debian:11-slim

LABEL version="0.1.0"
LABEL description="Kanaco Kana Converter Library"
LABEL maintainer="Takamichi Seki <takamichi.seki@elfincafe.com>"


RUN set -eux; \
    apt update; \
    apt -y full-upgrade; \
	apt -y install \
	  vim \
	  wget \
	  gpg \
	  doas \
	  g++ \
      gcc \
      libc6-dev \
      make \
      pkg-config \
	; \
	rm -rf /var/cache/apt/archives/*

ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.18.1

RUN set -eux; \
	arch="$(dpkg --print-architecture)"; arch="${arch##*-}"; \
	url=; \
	case "$arch" in \
		'amd64') \
			url='https://dl.google.com/go/go1.18.1.linux-amd64.tar.gz'; \
			sha256='b3b815f47ababac13810fc6021eb73d65478e0b2db4b09d348eefad9581a2334'; \
			;; \
		'arm64') \
			url='https://dl.google.com/go/go1.18.1.linux-arm64.tar.gz'; \
			sha256='56a91851c97fb4697077abbca38860f735c32b38993ff79b088dac46e4735633'; \
			;; \
		*) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; \
	esac; \
	build=; \
	if [ -z "$url" ]; then \
# https://github.com/golang/go/issues/38536#issuecomment-616897960
		build=1; \
		url='https://dl.google.com/go/go1.18.1.src.tar.gz'; \
		sha256='efd43e0f1402e083b73a03d444b7b6576bb4c539ac46208b63a916b69aca4088'; \
		echo >&2; \
		echo >&2 "warning: current architecture ($arch) does not have a compatible Go binary release; will be building from source"; \
		echo >&2; \
	fi; \
	\
	wget -O go.tgz.asc "$url.asc"; \
	wget -O go.tgz "$url" --progress=dot:giga; \
	echo "$sha256 *go.tgz" | sha256sum -c -; \
	\
# https://github.com/golang/go/issues/14739#issuecomment-324767697
	GNUPGHOME="$(mktemp -d)"; export GNUPGHOME; \
# https://www.google.com/linuxrepositories/
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC  EC91 7721 F63B D38B 4796'; \
# let's also fetch the specific subkey of that key explicitly that we expect "go.tgz.asc" to be signed by, just to make sure we definitely have it
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys '2F52 8D36 D67B 69ED F998  D857 78BD 6547 3CB3 BD13'; \
	gpg --batch --verify go.tgz.asc go.tgz; \
	gpgconf --kill all; \
	rm -rf "$GNUPGHOME" go.tgz.asc; \
	\
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	if [ -n "$build" ]; then \
		savedAptMark="$(apt-mark showmanual)"; \
		apt-get update; \
		apt-get install -y --no-install-recommends golang-go; \
		\
		( \
			cd /usr/local/go/src; \
# set GOROOT_BOOTSTRAP + GOHOST* such that we can build Go successfully
			export GOROOT_BOOTSTRAP="$(go env GOROOT)" GOHOSTOS="$GOOS" GOHOSTARCH="$GOARCH"; \
			./make.bash; \
		); \
		\
		apt-mark auto '.*' > /dev/null; \
		apt-mark manual $savedAptMark > /dev/null; \
		apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
		rm -rf /var/lib/apt/lists/*; \
		\
# remove a few intermediate / bootstrapping files the official binary release tarballs do not contain
		rm -rf \
			/usr/local/go/pkg/*/cmd \
			/usr/local/go/pkg/bootstrap \
			/usr/local/go/pkg/obj \
			/usr/local/go/pkg/tool/*/api \
			/usr/local/go/pkg/tool/*/go_bootstrap \
			/usr/local/go/src/cmd/dist/dist \
		; \
	fi; \
	\
	go version




RUN useradd -u 2000 -m -G staff,sudo zlkv
RUN echo "permit nopass :root" >> /etc/doas.conf
RUN echo "permit nopass :sudo" >> /etc/doas.conf
ENV GOROOT /usr/local/go
ENV GOPATH /var/opt
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/pkg" "$GOPATH/bin" && chown -R zlkv:staff "$GOPATH/pkg" "$GOPATH/bin"
WORKDIR $GOPATH/src
RUN touch ~/.profile
RUN echo 'export PS1="[\A]\u@kanaco:\W\\$ "' >> ~/.profile
USER zlkv
RUN touch ~/.profile
RUN echo 'export PS1="[\A]\u@kanaco:\W\\$ "' >> ~/.profile

