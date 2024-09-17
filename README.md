## Using agnostic-clickhouse-udf with clickhouse-local

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config-file ./examples/clickhouse/config.xml
```

### Load SQL UDFs (will be automated in the future)

```sql
CREATE OR REPLACE FUNCTION evm_hex_decode AS (s) -> unhex(substring(s, 3));
CREATE OR REPLACE FUNCTION evm_hex_encode AS (s) -> concatAssumeInjective('0x', lower(hex(s)));
CREATE OR REPLACE FUNCTION evm_hex_decode_uint256 AS (s) -> reinterpretAsUInt256(reverse(evm_hex_decode(trim(BOTH '"' FROM s))));
CREATE OR REPLACE FUNCTION evm_hex_encode_uint256 AS (i) -> concatAssumeInjective('0x', trim(LEADING '0' FROM lower(hex(i))));
```




