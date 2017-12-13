FROM alpine:3.6

ENV ALPINE_GLIBC_VERSION 2.26-r0
ENV HELM_VERSION 2.6.2

RUN apk update && apk --no-cache add ca-certificates openssl curl wget git bash \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
    && wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$ALPINE_GLIBC_VERSION/glibc-$ALPINE_GLIBC_VERSION.apk \
    && apk add glibc-$ALPINE_GLIBC_VERSION.apk

RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz -qO- | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64

RUN mkdir -p ~/.helm/plugins \
    && cd ~/.helm/plugins \
    && git clone https://github.com/app-registry/appr-helm-plugin.git registry \
    && helm registry --help

ADD draughtsman /

ENTRYPOINT ["/draughtsman"]
