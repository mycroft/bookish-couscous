FROM bookish-couscous-base:latest

RUN cd /go/src/gitlab.mkz.me/mycroft/bookish-couscous/client && \
    go get -v ./... && \
    go install gitlab.mkz.me/mycroft/bookish-couscous/client

ENTRYPOINT /bin/sh /go/src/gitlab.mkz.me/mycroft/bookish-couscous/sleeper.sh
