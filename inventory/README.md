# Inventory
This is a very small microservice to learn how to code using Go and Go-Kit.
It is written with clean architecture in mind and it implements logging and instrumentation using Prometheus.

## Building in Openshift
Follow what is described [here](/openshift/inventory/build/README.md)

## Building as a Docker Image
* **Install S2I** - https://github.com/openshift/source-to-image/releases/
* Clone *https://github.com/comolago/docker-images*
* Build both *go* and *alpine* Docker images (sources are under *s2i* directory)
* **Build** the application container image as follows:

```
cd shop
s2i build inventory registry.domain.local:5000/prod/s2i/go:1.11 registry.domain.local:5000/test/apps/inventory:1.0 --runtime-image=registry.domain.local:5000/prod/s2i/alpine:3.8
```

In this example s2i images are stored into *registry.domain.local:5000* registry beneath */prod/s2i* and the resulting image is tagged in such a way to be pushed on the same registry beneath *test/apps*

## Testing
To Test it, launch it and then issue curl commands such as:
```
curl -XPOST -d'{"id":"11", "Name": "Fedora Red" }' localhost:8080/items/add
```
