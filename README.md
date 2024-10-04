## Using agnostic-clickhouse-udf with clickhouse-local

### Build bundle

```sh
make bundle
```

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config ./examples/clickhouse-local-config.xml --path tmp/clickhouse
```

### Running clickhouse-server with Docker Compose

```sh
docker-compose up -d
```

