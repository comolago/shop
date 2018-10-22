# Golang Image
Edit the files to fit version your application you want to publish and the version of the Alpine Linux you want to use to embed it.

For example, to set it to embed version 0.8 of your application on Alpine Linuc 3.7
```
spec:
  output:
    to:
      kind: ImageStreamTag
      name: inventory:0.8
  source:
    dockerfile: |-
      FROM alpine:3.7
    images:
    - from: 
        kind: ImageStreamTag
        name: inventory-artifacts:0.8
  strategy:
    dockerStrategy:
      from: 
        kind: ImageStreamTag
        name: alpine-s2i:3.7
```
otherwise, if the defined version are the ones you want, simply leave it unmodified

Switch to the project you created, for example
```
oc project shop
```
Create the imagestream
```
oc apply -f imagestream.yaml  
```
Create The build configuration
```
oc apply -f buildconfig.yaml 
```
Build should automatically start

Wait until it completes
```
oc get builds
NAME       TYPE      FROM          STATUS     STARTED              DURATION
inventory-1             Docker    Dockerfile    Complete   35 seconds ago   7s
```
