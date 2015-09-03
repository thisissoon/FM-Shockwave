# Ubuntu just works
FROM alpine:3.2
MAINTAINER SOON_ <dorks@thisissoon.com>

## Environment Variables
ENV GOPATH /shockwave
ENV PATH $PATH:$GOPATH/bin

# OS Dependencies
RUN apk update && apk add go git ca-certificates alsa-lib-dev build-base bash && rm -rf /var/cache/apk/*

# Set working Directory
WORKDIR /shockwave

# GPM (Go Package Manager)
RUN git clone https://github.com/pote/gpm.git \
    && cd gpm \
    && git checkout v1.3.2 \
    && ./configure \
    && make install

# Install Dependencies
COPY ./Godeps /shockwave/Godeps
RUN gpm install

# Copy source code into the myleene src directory so Go can build the package
COPY . /shockwave/src/github.com/thisissoon/FM-Shockwave

# Set our final working dir to be where the source code lives
WORKDIR /shockwave/src/github.com/thisissoon/FM-Shockwave

# Install the go package
RUN go install ./...
