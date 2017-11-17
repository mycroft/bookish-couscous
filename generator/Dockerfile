FROM bookish-couscous-base:latest

RUN cd /go/src/gitlab.mkz.me/mycroft/bookish-couscous/generator && \
    go get -v ./... && \
    go install gitlab.mkz.me/mycroft/bookish-couscous/generator

ENTRYPOINT /bin/sh /go/src/gitlab.mkz.me/mycroft/bookish-couscous/sleeper.sh
