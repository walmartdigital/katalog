# katalog

Catalog all your kubernetes services in order to understand what is running and where is running

## Roles

* collector: read k8s and publish to a server
* server: receive information from collector

## Run

You can run the application in two differents roles: 1) Server role for listen
to all the collector signal and expose an interface to consult the calatog and
2) Collector in charge of capture every kubernetes event regarding objetct that
you wan to monitor

### Run The Server

```bash
go run src/main.go -role server
```

### Run The Collector

```bash
go run src/main.go -role collector
```

