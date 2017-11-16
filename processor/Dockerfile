FROM bookish-couscous-base:latest

RUN cd /go/src/gitlab.mkz.me/mycroft/bookish-couscous/processor && \
    go get -v ./... && \
    go install gitlab.mkz.me/mycroft/bookish-couscous/processor

ENTRYPOINT /go/bin/processor

