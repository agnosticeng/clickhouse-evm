## Using agnostic-clickhouse-udf with clickhouse-local

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config-file ./examples/clickhouse/config.xml
```

### Running clickhouse-local with the custom config

```sh
clickhouse local --config examples/clickhouse/config.xml --path tmp/clickhouse
```


