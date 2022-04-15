# OpenTelemetry Collector Demo

To run the demo application terminal commands

```shell
docker-compose up -d
```
or

```shell
make up
```

The demo exposes the following backends:

- Jaeger at http://0.0.0.0:16686
- Zipkin at http://0.0.0.0:9411
- Prometheus at http://0.0.0.0:9090 
