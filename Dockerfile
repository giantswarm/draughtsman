FROM alpine:3.6

RUN apk update && apk --no-cache add ca-certificates openssl curl wget && \
  update-ca-certificates && \
  wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
  wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.25-r0/glibc-2.25-r0.apk && \
  apk add glibc-2.25-r0.apk

ENV HELM_VERSION 2.6.2

RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz -qO- | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64

RUN mkdir -p ~/.helm/plugins/ \
    && wget https://github.com/app-registry/appr-helm-plugin/releases/download/v0.5.1/registry-helm-plugin.tar.gz -qO- | tar xzf - registry \
    && mv ./registry ~/.helm/plugins/ \
    && ~/.helm/plugins/registry/cnr.sh upgrade-plugin

ADD draughtsman /

ENTRYPOINT ["/draughtsman"]
