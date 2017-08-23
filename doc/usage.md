`kubectl create ns draughtsman`

`helm registry pull quay.io/giantswarm/draughtsman-chart@1.0.0-sha`  

```
Installation:
  V1:
    Secret:
      Draughtsman:
        SecretYaml: |
          service:
            deployer:
              environment: incluster-minikube
              eventer:
                github:
                  oauthtoken: XXX
                  organisation: giantswarm
                  projectlist: g8s-prometheus
              installer:
                helm:
                  organisation: giantswarm
                  password: XXX
                  username: user
              notifier:
                slack:
                  channel: alertmanager-test
            slack:
              token: XXX-XXX-XXX
      Registry:
        PullSecret:
          DockerConfigJSON: |-
            {
              "auths": {
                "quay.io": {
                  "auth": "XXX"
                }
              }
            }
```  


`cp values.yaml giantswarm_draughtsman-chart_1.0.0-sha/draughtsman-chart/values.yaml `  

`helm upgrade --install --reset-values draughtsman giantswarm_draughtsman-chart_1.0.0-sha/draughtsman-chart/`  

`kubectl get pods -n draughtsman`
```
NAME                           READY     STATUS    RESTARTS   AGE
draughtsman-1205598562-8494z   1/1       Running   0          30s
```  

`curl --request POST   --url https://api.github.com/repos/giantswarm/someproject/deployments  --header 'authorization: token XXX'   --header 'content-type: application/json'   --data '{  "ref": "XXX",  "environment": "incluster-minikube",     "auto_merge": false }'`  
