FROM golang

RUN apt clean all && apt update && apt upgrade -y
RUN apt install -y protobuf-compiler

# Base image to create all other images
# docker build -t bookish-couscous-base .

# Download most of dependencies to make sub dockerfile faster to generate
RUN go get github.com/golang/geo/s2 && \
    go get github.com/gocql/gocql && \
    go get google.golang.org/grpc && \
    go get github.com/Shopify/sarama && \
    go get github.com/golang/protobuf/proto && \
    go get github.com/garyburd/redigo/redis && \
    go get github.com/bsm/sarama-cluster && \
    go get google.golang.org/grpc/reflection && \
    go get -u github.com/golang/protobuf/protoc-gen-go

ENV PATTH="/opt/go/bin:${PATH}"

ADD . /go/src/gitlab.mkz.me/mycroft/bookish-couscous
RUN cd /go/src/gitlab.mkz.me/mycroft/bookish-couscous/common && go generate

