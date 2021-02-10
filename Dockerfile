FROM alpine:3.10

ENV HELM_VERSION 3.5.1
ENV APPR_PLUGIN_VERSION 0.7.0
ENV HELM_EXPERIMENTAL_OCI 1

# add application user
RUN addgroup -S draughtsman && adduser -S -g draughtsman draughtsman

# dependencies
RUN set -x \
    && apk update && apk --no-cache add ca-certificates openssl curl bash zlib

# install helm
RUN set -x \
    && curl -s https://get.helm.sh/helm-v$HELM_VERSION-linux-amd64.tar.gz | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64

# install helm appr (registry) plugin
RUN helm plugin install https://github.com/app-registry/quay-helmv3-plugin

USER draughtsman

ADD draughtsman /home/draughtsman/

# setup default catalog repo
RUN helm repo add default-catalog https://giantswarm.github.io/default-catalog/

RUN helm repo update

WORKDIR /home/draughtsman

ENTRYPOINT ["/home/draughtsman/draughtsman"]
