## Using agnostic-clickhouse-udf with clickhouse-local

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config-file ./examples/clickhouse/config.xml
```

### Load SQL UDFs (will be automated in the future)

```sql
CREATE FUNCTION evm_hex_decode AS (s) -> unhex(substring(s, 3));
CREATE FUNCTION evm_hex_encode AS (s) -> concatAssumeInjective('0x', lower(hex(s)));
```


