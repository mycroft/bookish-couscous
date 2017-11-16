FROM bookish-couscous-base:latest

RUN cd /go/src/gitlab.mkz.me/mycroft/bookish-couscous/fo && \
    go get -v ./... && \
    go install gitlab.mkz.me/mycroft/bookish-couscous/fo

ENTRYPOINT /go/bin/fo