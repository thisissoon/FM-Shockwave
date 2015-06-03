# Ubuntu just works
FROM ubuntu:14.04
MAINTAINER SOON_ <dorks@thisissoon.com>

## Environment Variables
ENV GOPATH /shockwave
ENV PATH $PATH:$GOPATH/bin

# OS Dependencies
RUN apt-get update -y && apt-get install --no-install-recommends -y -q \
        build-essential \
        software-properties-common \
        libasound2-dev \
        git \
    && apt-get clean \
    && apt-get autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

# Add Go PPA
RUN apt-add-repository -y ppa:evarlast/golang1.4

# Install Go
RUN apt-get update && apt-get install -y golang

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
