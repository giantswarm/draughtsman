FROM alpine:3.6

#ENV ALPINE_GLIBC_VERSION 2.26-r0
ENV HELM_VERSION 2.6.2
ENV APPR_PLUGIN_VERSION 0.7.0

RUN set -x \
    && apk update && apk --no-cache add ca-certificates openssl curl bash zlib

# RUN set -x \
#     && curl -s https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub -o /etc/apk/keys/sgerrand.rsa.pub \
#     && curl -s -L https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$ALPINE_GLIBC_VERSION/glibc-$ALPINE_GLIBC_VERSION.apk -o ./glibc-$ALPINE_GLIBC_VERSION.apk \
#     && pwd && ls -la \
#     && apk add ./glibc-$ALPINE_GLIBC_VERSION.apk \
#     && rm ./glibc-$ALPINE_GLIBC_VERSION.apk


RUN set -x \
    && curl -s https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64


RUN set -x \
    && mkdir -p ~/.helm/plugins \
    && curl -L -s https://github.com/app-registry/appr-helm-plugin/releases/download/v$APPR_PLUGIN_VERSION/helm-registry_linux.tar.gz | tar xvzf - registry \
    && mv ./registry ~/.helm/plugins/registry \
    && ~/.helm/plugins/registry/cnr.sh upgrade-plugin \
    && helm registry --help >> /dev/null

ADD draughtsman /

ENTRYPOINT ["/draughtsman"]
