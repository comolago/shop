# Golang Image
Edit the files to fit version of the GO S2I image you want to use.

For example, to set it to build with Golang 1.7 inventory version 1.8, modify 
```
spec:
  output:
    to:
      kind: ImageStreamTag
      name: inventory-artifacts:0.8
  strategy:
    sourceStrategy:
      from:
        kind: ImageStreamTag
        name: go-s2i:1.7
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
Build it
```
oc start-build inventory-artifacts
```
Wait until it completes
```
oc get builds
NAME                    TYPE      FROM          STATUS     STARTED         DURATION
inventory-artifacts-1   Source    Git@60815d5   Complete   4 minutes ago   21s
```
