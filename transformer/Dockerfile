############################
# STEP 1 build executable binary
############################


FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/package/app/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v ./... && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/main

############################
# STEP 1 build executable binary
############################

FROM python:3.7.7-buster AS extracting
LABEL maintainer="yossefaz@protonmail.com"
RUN apt-get update && \
    apt-get install -y --no-install-recommends git autoconf libtool swig texinfo build-essential gcc python3-libxml2 && \
    LIBXML2VER=2.9.9 && \
    mkdir libxmlInstall && cd libxmlInstall && \
    wget ftp://xmlsoft.org/libxml2/libxml2-$LIBXML2VER.tar.gz && \
    tar xf libxml2-$LIBXML2VER.tar.gz && \
    cd libxml2-$LIBXML2VER/ && \
    ./configure && \
    make && \
    make install && \
    cd /libxmlInstall && \
    rm -rf gg libxml2-$LIBXML2VER.tar.gz libxml2-$LIBXML2VER
WORKDIR /app
RUN git clone git://git.sv.gnu.org/libredwg.git && \
     cd libredwg && \
     sh autogen.sh && \
     ./configure --enable-release --prefix=/usr/local && \
     make && \
     mkdir install && \
     make install DESTDIR="$PWD/install" && \
     make check DESTDIR="$PWD/install"

############################
# STEP 2 copy the executable binary to thinest image
############################

FROM debian:stable-slim
COPY --from=extracting /app/libredwg/install/usr/local/bin/* /usr/local/bin/
COPY --from=extracting /app/libredwg/install/usr/local/include/* /usr/local/include/
COPY --from=extracting /app/libredwg/install/usr/local/lib/* /usr/local/lib/
RUN ldconfig


WORKDIR /app

RUN mkdir -p shared/logs config && \
    chmod 777 shared
COPY --from=builder /go/bin/main .
COPY ./config/*.json ./config/
# Run the hello binary.
CMD ["./main"]



