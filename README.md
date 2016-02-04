# kompose

`kompose` is a fork of [libcompose](https://github.com/docker/libcompose) which is a Go library for [Docker Compose](http://docs.docker.com/compose).
`kompose` adds [Kubernetes](http://kubernetes.io) support. It takes a Docker Compose file and translates it into Kubernetes objects, it then submits those objects to a Kubernetes endpoint.

kompose is a convenience tool to go from local Docker development to managing your application with Kubernetes. We don't assume that the transformation from docker compose format to Kubernetes API objects will be perfect, but it helps tremendously to start _Kubernetizing_ your application.

## Download

Grab the latest [release](https://github.com/skippbox/kompose/releases)

For example on OSX:

```bash
$ sudo curl -L https://github.com/skippbox/kompose/releases/download/v0.0.1/kompose_darwin-amd64 > /usr/local/bin/kompose
$ sudo chmod +x /usr/local/bin/kompose
```

## Usage

You need a Docker Compose file handy. There is a sample one in the `samples/` directory for testing.
You will convert the compose file to K8s objects with `kompose k8s convert`

```bash
$ cd samples/
$ ls
docker-compose.yml
$ kubectl get rc
CONTROLLER   CONTAINER(S)   IMAGE(S)   SELECTOR   REPLICAS   AGE
$ kompose k8s convert
```

Check that the replication controllers and services have been created.
The .yaml file will be in the same directory.

```bash
$ kubectl get rc
CONTROLLER   CONTAINER(S)   IMAGE(S)                SELECTOR        REPLICAS   AGE
redis        redis          redis:3.0               service=redis   1          2s
web          web            tuna/docker-counter23   service=web     1          2s
$ ls
docker-compose.yml	redis-rc.yaml		redis-svc.yaml		web-rc.yaml
```

kompose also allows you to list the replication controllers and services with the `ps` subcommand.
You can delete them with the `delete` subcommand.

```bash
$ kompose k8s ps --rc
Name           Containers     Images                        Replicas  Selectors           
redis          redis          redis:3.0                     1         service=redis       
web            web            tuna/docker-counter23         1         service=web         
$ kompose k8s ps --svc
Name                Cluster IP          Ports               Selectors           
redis               10.0.20.194         TCP(6379)           service=redis       
$ kompose k8s delete --rc --name web
$ kompose k8s ps --rc
Name           Containers     Images                        Replicas  Selectors           
redis          redis          redis:3.0                     1         service=redis       
```

And finally you can scale a replication controller with `scale`.

```bash
$ kompose k8s scale --scale 3 --rc redis
Scaling redis to: 3
$ kompose k8s ps --rc
Name           Containers     Images                        Replicas  Selectors           
redis          redis          redis:3.0                     3         service=redis       
$ kubectl get rc
CONTROLLER   CONTAINER(S)   IMAGE(S)    SELECTOR        REPLICAS   AGE
redis        redis          redis:3.0   service=redis   3          45s
```

Note that you can of course manage the services and replication controllers that have been created with `kubectl`.
The command of kompose have been extended to match the `docker-compose` commands.

## Alternate formats

The default `kompose` transformation will generate replication controllers and services. You can alternatively generate [Deployment](https://github.com/kubernetes/kubernetes/blob/release-1.1/docs/user-guide/managing-deployments.md) objects or [Helm](https://github.com/helm/helm) charts.

```bash
$ kompose k8s convert --d
$ ls
$ tree
.
├── docker-compose.yml
├── redis-deployment.yaml
├── redis-rc.yaml
├── redis-svc.yaml
├── web-deployment.yaml
└── web-rc.yaml
```

The `*deployment.yaml` files contain the Deployments objects

```bash
$ kompose k8s convert --c
$ tree docker-compose/
docker-compose/
├── Chart.yaml
├── README.md
└── manifests
    ├── redis-rc.yaml
    ├── redis-svc.yaml
    └── web-rc.yaml
```

The chart structure is aimed at providing a skeleton for building your Helm charts.

## Building

You need either [Docker](http://github.com/docker/docker) and `make`,
or `go` in order to build libcompose. To simplify this a Vagrantfile is provided.

### Building with `Vagrant`

After cloning the repo, `vagrant up` and use `make` inside the machine

```bash
$ git clone https://github.com/skippbox/kompose.git
$ cd kompose
$ vagrant up
$ vagrant ssh
$ cd /vagrant
$ make binary
```

The binaries will be in the `bundles/` directory

```bash
$ ls bundles/
kompose_darwin-386		kompose_linux-386		kompose_linux-arm		kompose_windows-amd64.exe
kompose_darwin-amd64		kompose_linux-amd64		kompose_windows-386.exe
```

### Building with `docker`

You need Docker and ``make`` and then run the ``binary`` target. This
will create binary for all platform in the `bundles` folder. 

```bash
$ make binary
docker build -t "libcompose-dev:refactor-makefile" .
# […]
---> Making bundle: binary (in .)
Number of parallel builds: 4

-->      darwin/386: github.com/docker/libcompose/cli/main
-->    darwin/amd64: github.com/docker/libcompose/cli/main
-->       linux/386: github.com/docker/libcompose/cli/main
-->     linux/amd64: github.com/docker/libcompose/cli/main
-->       linux/arm: github.com/docker/libcompose/cli/main
-->     windows/386: github.com/docker/libcompose/cli/main
-->   windows/amd64: github.com/docker/libcompose/cli/main

```

### Building with `go`

- You need `go` v1.5
- You need to set export `GO15VENDOREXPERIMENT=1` environment variable
- If your working copy is not in your `GOPATH`, you need to set it
accordingly.

```bash
$ go generate
# Generate some stuff
$ go build -o libcompose ./cli/main
```

## Contributing and Issues

`kompose` is a work in progress, we will see how far it takes us. We welcome any pull request to make it even better.
If you find any issues, please [file it](https://github.com/skippbox/kompose/issues).
