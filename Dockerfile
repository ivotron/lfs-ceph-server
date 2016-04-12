FROM golang:1.6

MAINTAINER Ivo Jimenez <ivo.jimenez@gmail.com>

RUN apt-get update && \
    apt-get install -y --no-install-recommends  librados-dev && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

ADD ceph.conf /root

EXPOSE 8080
EXPOSE 9999

CMD ["go-wrapper", "run"]
