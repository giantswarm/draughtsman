## Introduction
This document is going to explain the basic usage and setup of draughtsman with a kubernetes cluster.

## Preparation
The following steps have to be completed before starting with the actual installation of draughtsman to a k8s cluster.
#### Setting up Minikube (optional)
You can skip this step if you have a running k8s cluster and want to run draughtsman on it.
Follow the installation instructions provided by minikube for your system. The output of `minikube status` should be:
```
> minikube status
minikube:
localkube:
```
if the installation was successfull. </br>
Start up kubernetes with the command:
`minikube start --cpus 4 --memory 4096 --kubernetes-version v1.7.0` (Alter the memory and cpu requirements if you do not have the resources necessary).
#### Setting up Helm
Both helm and the registry helm plugin are necessary to run draughtsman. Follow the helm installation instructions [here](someurl) in order to install helm and then apply the registry instructions [here](someotherrul) to install the registry plugin on top.</br>
Running `helm registry -h` should return the registry plugin help info, if the installation was successful.
```
usage: appr [-h]
            {pull,run-server,show,inspect,list,delete-package,helm,version,logout,deploy,plugins,push,login,config,channel,jsonnet}
            ...
```
Run `helm init` in order to install tiller into your cluster. Check the state of the tiller deployment (on top of the return message from helm) by running `kubectl describe deployment tiller-deploy -n=kube-system`. The deployment should be satisfied if tiller was installed correctly.

#### Configuring Draughtsman
Create a draughtsman namespace by running `kubectl create ns draughtsman`. Confirm that the namespace exists by running `kubectl get ns`.

Download the draughtsman chart with the helm registry `helm registry pull quay.io/giantswarm/draughtsman-chart@1.0.0-{sha}`. This should create a folder with the chart in your current directory.

Create a new file names `values.yaml` with the following content:
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
                  organisation: XXX
                  projectlist: XXX,XXX
              installer:
                helm:
                  organisation: XXX
                  password: XXX
                  username: user
              notifier:
                slack:
                  channel: XXX-XX
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
Several values in this file need to be adjusted to your environment:
* `github/oauthtoken` Your oauthtoken which has deploy rights in your github organisation.
* `github/organisation` Your github organisation.
* `github/projectlist` A comma seperated list of projects in your organisation which are covered by the oauthtoken.
* `helm/organisation` Your organisation on your chart repository (quay.io for example).
* `helm/user` A user authorized for your chart repository.
* `helm/password` The users password.
* `slack/channel` The slack channel you want draughtsman to report to.
* `slack/token` A slack api token for your slack.
* `Registry/PullSecret` The pull secret is needed if your docker images are in a private repository.

Replace the empty values file in the downloaded chart with the new one:
`cp values.yaml giantswarm_draughtsman-chart_1.0.0-{sha}/draughtsman-chart/values.yaml `

## Installation
The actual installation, after you've configured draughtsman correctly, is simply running the following command:</br>
`helm upgrade --install --reset-values draughtsman giantswarm_draughtsman-chart_1.0.0-sha/draughtsman-chart/`  


This should result into a running pod of draughtsman, executing `kubectl get pods -n draughtsman -n draughtsman` should return something similar to:
```
NAME                           READY     STATUS    RESTARTS   AGE
draughtsman-1205598562-8494z   1/1       Running   0          30s
```  

Check the logs of draughtsman with `kubectl logs draughtsman-1205598562-8494z`(your podname of course) which should look something like this if everything went well after startup (cleaned up for readability):
```
{... "debug":"creating in-cluster config", ...}
{... "debug":"checking connection to Kubernetes", ...}
{... "debug":"checking connection to Kubernetes", ...}
{... "debug":"logging into registry","registry":"XXX,... ,"username":"XXX"}
{... "debug":"running helm command","name":"login", ...}
{... "debug":"ran helm command","name":"login","stderr":"","stdout":" \u003e\u003e\u003e Login succeeded\n", ...}
{... "debug":"checking connection to Slack", ...}
{... "debug":"starting deployer", ...}
{... "debug":"starting polling for github deployment events","interval":"1m0s", ...}
{... "debug":"fetching deployments","project":"XXX", ...}
```

### Debugging errors
Draughtsman will panic if a configuration is missing, all entries in the `values.yaml` are mandatory.

Draughtsman will try to run `helm registry login` with the supplied credentials on startup. The logs should therefor contain something along the lines of `login successfull`. Draughtsman will not panic if the login wasn't successfull, but it will not operate correctly. Errormessages that include `install registry,...` indicate that tiller is not running in your cluster.
## Configuring applications
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: draughtsman-values-configmap
  namespace: draughtsman
data:
  values: |
    Installation:
      V1:
        Monitoring:
          Alertmanager:
            Address: "abc"
            Host: "abc"
          Prometheus:
            Address: "abc"
            Host: "abc"
            ClusterLabel: 'incluster'
            RetentionPeriod: "336h"
```
```
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: draughtsman-values-secret
  namespace: draughtsman
data:
  values: BASE64
```
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
                  organisation: XXX
                  projectlist: XXX,XXX
              installer:
                helm:
                  organisation: XXX
                  password: XXX
                  username: user
              notifier:
                slack:
                  channel: XXX-XX
            slack:
              token: XXX-XXX-XXX
      Prometheus:
        Nginx:
          Auth: XXX:XXX
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
## Deployments
Github deployments events are created by interacting with the github API directly. Run the following curl command in order to create a deployment event manually:
```
curl --request POST   --url https://api.github.com/repos/yourorg/yourproject/deployments  --header 'authorization: token XXX'   --header 'content-type: application/json'   --data '{  "ref": "XXX",  "environment": "incluster-minikube",     "auto_merge": false }' ```

Note that the project has to be in the projectlist of the `values.yaml` from above. The authorization token has to be a github oauthtoken with deploy rights and the `ref` has to be the commit SHA.

```
{... "debug":"fetching deployments","project":"XXX", ...}
{... "debug":"found new deployment events","project":"XXX","time":"17-08-31 08:40:55.047"}
{... "debug":"posting deployment status","id":XXX,"project":"XXX","state":"pending", ...}
{... "debug":"installing chart","name":"XXX","sha":"XXX", ...}
{... "debug":"running helm command","name":"pull", ...}

```
## Metrics and alerting
Draughtsman exposes several metrics, which can be used to alert on draughtsmans behavior:
```
draughtsman_configmap_configurer_request_duration_milliseconds
draughtsman_configmap_configurer_request_total

draughtsman_github_eventer_github_deployment_duration_milliseconds
draughtsman_github_eventer_github_deployment_response_code
draughtsman_github_eventer_github_deployment_status_duration_milliseconds
draughtsman_github_eventer_github_deployment_status_response_code
draughtsman_github_eventer_rate_limit_limit
draughtsman_github_eventer_rate_limit_remaining

draughtsman_helm_installer_helm_command_duration_milliseconds
draughtsman_helm_installer_helm_command_total

draughtsman_slack_notifier_request_duration_milliseconds
draughtsman_slack_notifier_request_total

```
