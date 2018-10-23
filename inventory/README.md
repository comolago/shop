=Inventory
This is a very small microservice to learn how to code using Go and Go-Kit.
It is written with clean architecture in mind and it implements logging and instrumentation using Prometheus.

To Test it, launch it and then issue curl commands such as:
```
curl -XPOST -d'{"id":"11", "Name": "Fedora Red" }' localhost:8080/items/add
```
