# draughtsman

[![CircleCI](https://circleci.com/gh/giantswarm/draughtsman.svg?&style=shield)](https://circleci.com/gh/giantswarm/draughtsman) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/draughtsman/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/draughtsman)

`draughtsman` is a deployment agent for Kubernetes clusters.

It is designed to be used in several Kubernetes clusters to deploy and manage applications running with different configurations. 

# Kubernetes Configuration

`draughtsman` runs in its own namespace, named `draughtsman`.

For configuration of `draughtsman` itself, a Secret named `draughtsman` is required, inside the `draughtsman` namespace.

For example:
```
apiVersion: v1
kind: Secret
metadata:
  name: draughtsman
  namespace: draughtsman
  labels:
    app: draughtsman
type: Opaque
data:
  secret.yaml: <SECRET-CONFIGURATION-BASE64-ENCODED-HERE>
```

with the secret configuration (for example) as follows:

```
service:
    slack:
        token: ...
    deployer:
        environment: ...
        eventer:
            github:
                oauthtoken: ...
                organisation: ...
                projectlist: ...
        installer:
            helm:
                username: ...
                password: ...
                organisation: ...
        notifier:
            slack:
                channel: ...
```

This file needs to be updated, the file contents base64 encoded, and then inserted into the secret.

A second ConfigMap is also necessary, if using the ConfigMap Configurer (the default).
By default, this configmap is named `draughtsman-values`, and is in the `draughtsman` namespace.

For example:
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: draughtsman-values
  namespace: draughtsman
data:
  values: |
    Installation:
      V1:
        GiantSwarm:
          API:
            Address:
              Scheme: "https"
              Host: "api-test.giantswarm.io
```

All data under the `values` key (by default), is passed verbatim to Helm, to provide values for chart Installations.
