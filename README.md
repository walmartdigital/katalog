# Katalog

Catalog all your kubernetes services in order to understand what is running and where is running

## Roles

* collector: read k8s and publish to a server
* server: receive information from collector and expose it as a REST API

## Run

You can run the application in two differents roles: 1) Server role for listen
to all the collector signal and expose an interface to consult the catalog and
2) Collector in charge of capture every kubernetes event regarding object that
you wan to monitor

### Run The Server

```bash
go run src/main.go -role server
```

### Run The Collector (developer mode)

```bash
go run src/main.go -kubeconfig
```

> When deploying on K8s cluster, omit the `-kubeconfig` flag

### Env Variables

- **PUBLISHER:** How to publish events. Values can be http or kafka (default http)
- **LOG_LEVEL:** Log level. Values can be DEBUG, WARN, INFO or ERROR (default ERROR)
- **HTTP_URL:** Url to use with http publisher
- **KAFKA_URL:** Url to use with kafka publisher
- **KAFKA_TOPIC_PREFIX:** topic prefix to use on kafka publisher. Default ```_katalog.artifcat.```


### Run local environment

A development environment is avalaible using skaffold.

- brew install minikube
- minikube start
- skaffold dev

Alternative with kind:

- brew install kind
- kind create cluster
- skaffold dev
