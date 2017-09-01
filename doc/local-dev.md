#### Goal
The following steps are going to set up a dev environment with minikube. This is meant to allow you to test your changes to draughtsman before they are committed or merged.

*The usage guide is more detailed about actually running draughtsman. This guide is only meant to set up a local environment to develop on draughtsman!*

#### Preparation

The following steps are meant to prepare your minikube and the environment `draughtsman` is going to need. You should have `minikube` and docker installed on your machine beforehand.  
Start up your minikube, here is an example configuration (note that `k8s v1.7.3` caused some errors in our setups):  
`minikube start --cpus 4 --memory 4096 --kubernetes-version v1.7.0`  
Your kubectl should be configured to use your minikube now and kubernetes should be available at `https://192.168.99.100:8443`.  
Next install helm (see installation guides [here](https://github.com/kubernetes/helm/blob/master/docs/install.md)) and run `helm init` against your minikube. Check if tiller is successfully deployed with:  
 `kubectl describe deployment tiller-deploy -n=kube-system`.

A `secret.yaml` has to be prepared next, which contains all information for draughtsman to connect to your minikube and other entities:
```
service:
  kubernetes:
    incluster: false
    address: https://192.168.99.100:8443
    tls:
      cafile: /minikube/apiserver.crt
      crtfile: /minikube/apiserver.crt
      keyfile: /minikube/apiserver.key
  deployer:
    environment: minikube
    eventer:
      github:
        oauthtoken: your_token
        organisation: your_org
        projectlist: your_project, your_other_project
    installer:
      helm:
        organisation: your_quay_org
        password: your_quay_password
        username: your_quay_username
    notifier:
      slack:
        channel: your_channel
  slack:
    token: your_slack_token
```
Note that *all* of these values are mandatory! `Draughtsman` will not work if one or several of these values are missing!

You now have to prepare your cluster itself now, a secret called `draughtsman-values-secret` and a configmap called `draughtsman-values-configmap` are necessary for draughtsman to run. Make sure a namespace called `draughtsman` exists.  
 The configmap can look something like this:
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
        YOURORG:
          YOURAPP:
            Address:
              Scheme: "something"
              Host: "something"
```
Note that these are values match to the placeholders in your app ( see draughtsman-usage ).

The secret should look like this:
```
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: draughtsman-values-secret
  namespace: draughtsman
data:
  values: SOMEBASE64
```
Where the base64 value should match this pattern (see draughtsman-usage):
```
Installation:
  V1:
    Secret:
      APP:
        Nginx:
          Auth: SOMETHING
      APP2:
        Nginx:
          Auth: SOMETHING
```
Now you can add both to your cluster with `kubectl apply -f filename`.

#### Dev-cycle
Once you've made a change in draughtsman first use `go build` to create the draughtsman binary.  
Next up create the docker image: `docker build . --tag draughtsman-dev`  
And then go ahead and start the container with some configs:  
`docker run -v ~/.minikube:/minikube -v /your/folder:/var/run/draughtsmanconfig  draughtsman-dev daemon --config.dirs=/var/run/draughtsmanconfig --config.files=secret`  
Note that `/your/folder` is the folder with your `secret.yaml` from the preparation part and that `~/.minikube` is assumed to be your minikube folder locally.  

You should now see that draughtsman is polling for new deployment events in `your_project`. You can use curl in order to create new deployment events:  
`curl --request POST   --url https://api.github.com/repos/your_org/your_project/deployments  --header 'authorization: token your_token'   --header 'content-type: application/json'   --data '{  "ref": "your_commit",  "environment": "minikube",     "auto_merge": false }'`  

This should be all you need in order to observe your changes locally.

