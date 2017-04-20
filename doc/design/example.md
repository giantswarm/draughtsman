### build the container

this would be part of the `architect` build

```
$ docker build -t quay.io/josephsalisbury/helloworld:test .
Sending build context to Docker daemon 10.72 MB
Step 1/5 : FROM busybox:ubuntu-14.04
 ---> d16744963217
Step 2/5 : ADD ./helloworld /usr/bin/
 ---> Using cache
 ---> 5e6f54c5aadc
Step 3/5 : ADD content /content
 ---> Using cache
 ---> 4c4e6d1b368b
Step 4/5 : EXPOSE 8080
 ---> Using cache
 ---> e4fba2fb5a9c
Step 5/5 : ENTRYPOINT helloworld
 ---> Using cache
 ---> fc60039a28c8
Successfully built fc60039a28c8
```
```
$ docker images | grep quay | grep helloworld
quay.io/josephsalisbury/helloworld                     test                                                      fc60039a28c8        3 minutes ago       11.7 MB
```

### push the container

this would be part of the `architect` build - i want to build and push all images

```
$ docker push quay.io/josephsalisbury/helloworld:test
The push refers to a repository [quay.io/josephsalisbury/helloworld]
df573003d969: Pushed
1aa9eeb6a585: Pushed
5f70bf18a086: Pushed
5dbcf0efe4f2: Pushed
test: digest: sha256:63ce20d0c374c89082e8da1f85009d74b7616064ce3096fff85543a14c6ac893 size: 4792
```