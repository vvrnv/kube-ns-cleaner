# kube-ns-cleaner

![image](https://user-images.githubusercontent.com/40491079/93712746-651c3100-fb60-11ea-89e9-4c207f7db7ef.png)

kube-ns-cleaner is a kubernetes client application for deleting unused kubernetes namespaces and scaling their resources

## About The Project

When managing bare-metal k8s clusters, you may encounter the problem of limited resources. An example of a large number of users using the same cluster, or CI/CD for some reason did not delete the namespace after testing. This application deletes a namespaces or scaling to zero deployments, stateful sets, and removing daemonsets.

## Running from source code

```sh
go run main.go


```

## Installation

### Helm

- [] todo

## Download binary

[release page link](https://github.com/vvrnv/kube-ns-cleaner/releases)
