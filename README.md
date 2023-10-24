# Sending Prometheus metrics to DataDog

Sending most gauge metrics to DataDog

## How to build
```
go build -o prom2dd -tags musl src/main.go
```

## Build image
```
make container
```