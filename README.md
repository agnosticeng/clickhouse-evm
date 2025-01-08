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

### Installing a bundle in a running clickhouse-server container

```sh
su - clickhouse -c "wget -qO- https://github.com/agnosticeng/agnostic-clickhouse-udf/releases/download/v0.0.5/agnostic-clickhouse-udf_0.0.5_linux_amd64_v3.tar.gz | tar xvz -C /"
for f in /var/lib/clickhouse/user_defined/*.sql; do clickhouse client --queries-file $f; done
```
