### evm_hex_encode_int

Encode the input integer as hexadecimal with 0x prefix.

**Syntax**

```sql
evm_hex_encode_int(str)
```

**Parameters**

- `int` - Any string. [Int|UInt](https://clickhouse.com/docs/en/sql-reference/data-types/int-uint)

**Returned value**

- Returns the 0x-prefixed hexadecimal encoding of the input integer.

**Example**

Query:

```sql
select evm_hex_encode_int(123)
```

Result:

| evm_hex_encode_int(123) |
|:-|
| 0x7b |