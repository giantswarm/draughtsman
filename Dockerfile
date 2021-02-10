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
RUN set -x \
    && mkdir -p /home/draughtsman/.helm/plugins \
    && curl -L -s https://github.com/app-registry/appr-helm-plugin/releases/download/v$APPR_PLUGIN_VERSION/helm-registry_linux.tar.gz | tar xvzf - registry \
    && mv ./registry /home/draughtsman/.helm/plugins/registry \
    && chown -R draughtsman:draughtsman /home/draughtsman/.helm

# setup default catalog repo
RUN helm repo add default-catalog https://giantswarm.github.io/default-catalog/ && helm repo update

USER draughtsman

ADD draughtsman /home/draughtsman/

RUN cd /home/draughtsman/.helm/plugins/registry \
    && ./cnr.sh upgrade-plugin \
    && helm registry --help > /dev/null

WORKDIR /home/draughtsman

ENTRYPOINT ["/home/draughtsman/draughtsman"]
