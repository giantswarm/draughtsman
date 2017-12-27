FROM alpine:3.6

ENV HELM_VERSION 2.6.2
ENV APPR_PLUGIN_VERSION 0.7.0

# add application user
RUN addgroup -S app && adduser -S -g draughtsman draughtsman

# dependencies
RUN set -x \
    && apk update && apk --no-cache add ca-certificates openssl curl bash zlib

# install helm
RUN set -x \
    && curl -s https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64

# install helm appr (registry) plugin
RUN set -x \
    && mkdir -p ~/.helm/plugins \
    && curl -L -s https://github.com/app-registry/appr-helm-plugin/releases/download/v$APPR_PLUGIN_VERSION/helm-registry_linux.tar.gz | tar xvzf - registry \
    && mv ./registry ~/.helm/plugins/registry \
    && ~/.helm/plugins/registry/cnr.sh upgrade-plugin \
    && helm registry --help >> /dev/null

ADD draughtsman /

USER draughtsman

ENTRYPOINT ["/draughtsman"]
