FROM ubuntu:xenial

RUN apt-get -y update \
    && apt-get -y install \
    wget

RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v2.6.2-linux-amd64.tar.gz -qO- | tar xzf - linux-amd64/helm \
    && chmod +x ./linux-amd64/helm \
    && mv ./linux-amd64/helm /bin/helm \
    && rm -rf ./linux-amd64

RUN mkdir -p ~/.helm/plugins/ \
    && wget https://github.com/app-registry/appr-helm-plugin/releases/download/v0.5.1/registry-helm-plugin.tar.gz -qO- | tar xzf - registry \
    && mv ./registry ~/.helm/plugins/

ADD draughtsman /

ENTRYPOINT ["/draughtsman"]
