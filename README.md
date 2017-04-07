mystack-logger
==============
[![Build Status](https://travis-ci.org/topfreegames/mystack-logger.svg?branch=master)](https://travis-ci.org/topfreegames/mystack-logger)

The logger component of mystack

### About
This is the mystack logger component, it will aggregate logs from all apps deployed by mystack on Kubernetes and expose them through an API.

### Dependencies
* Go 1.8
* Redis
* NSQ

### Building
```
make build
```

### Running
```
make deps
make run
```

### Build a docker image
```
make build-docker
```

### Installing into kubernetes
```
kubectl create namespace mystack
kubectl apply -f ./manifests
```
