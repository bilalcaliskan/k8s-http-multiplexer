# K8s Http Multiplexer
[![CI](https://github.com/bilalcaliskan/k8s-http-multiplexer/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/k8s-http-multiplexer/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/bilalcaliskan/k8s-http-multiplexer)](https://hub.docker.com/r/bilalcaliskan/k8s-http-multiplexer/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/k8s-http-multiplexer)](https://goreportcard.com/report/github.com/bilalcaliskan/k8s-http-multiplexer)
[![codecov](https://codecov.io/gh/bilalcaliskan/k8s-http-multiplexer/branch/master/graph/badge.svg)](https://codecov.io/gh/bilalcaliskan/k8s-http-multiplexer)
[![Release](https://img.shields.io/github/release/bilalcaliskan/k8s-http-multiplexer.svg)](https://github.com/bilalcaliskan/k8s-http-multiplexer/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/k8s-http-multiplexer)](https://github.com/bilalcaliskan/k8s-http-multiplexer)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a project for multiplexing HTTP requests inside a Kubernetes cluster. When you need to send a HTTP request to each container in
a cluster cluster with a single HTTP request, k8s-http-multiplexer is what you need exactly.

> **k8s-http-multiplexer should be running inside a target Kubernetes cluster to properly operate**

## Installation
K8s-http-multiplexer can be deployed as Kubernetes deployment or standalone installation. You can use [sample deployment file](deployment/deployment_all.yaml) to deploy your Kubernetes cluster.
Before make deployment, you should deploy the [sample configmap](deployment/configmap.yaml) to the cluster. k8s-http-multiplexer reads that
configmap to take proper actions on each HTTP request.
```shell
$ kubectl create -f deployment/configmap.yaml
$ kubectl create -f deployment/deployment_all.yaml
```

## Configuration
K8s-http-multiplexer can be customized with several command line arguments at the app level, and a configmap at the business level.
Here is the list of arguments you can pass:
```
--kubeConfigPath        Kube config file path to access cluster. Required while running out of Kubernetes cluster.
--configFilePath        Path of the config file of k8s-http-multiplexer to read, defaults to /opt/config/config.yaml
--inCluster             Boolean variable if k8s-http-multiplexer is running inside k8s cluster or not, required for
                        debugging purpose. Defaults to true.
```

You can inspect [sample config file](config/sample.yaml) and [sample configmap object](deployment/configmap.yaml).

## Development
This project requires below tools while developing:
- [Golang 1.16](https://golang.org/doc/go1.16)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)

## License
Apache License 2.0

## How k8s-http-multiplexer handles authentication/authorization with kube-apiserver?

k8s-http-multiplexer uses [client-go](https://github.com/kubernetes/client-go) to interact
with `kube-apiserver`. [client-go](https://github.com/kubernetes/client-go) uses the [service account token](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/)
mounted inside the Pod at the `/var/run/secrets/kubernetes.io/serviceaccount` path while initializing the client.

If you have RBAC enabled on your cluster, when you applied the sample deployment file [config/sample.yaml](deployment/deployment_all.yaml),
it will create required serviceaccount, role and rolebinding and then use that serviceaccount to be used
by our k8s-http-multiplexer pods.

If RBAC is not enabled on your cluster, please follow [that documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) to enable it.
