### evm_hex_encode

Encode the input string as hexadecimal with 0x prefix.

**Syntax**

```sql
evm_hex_encode(str)
```

**Parameters**

- `str` - Any string. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns the 0x-prefixed hexadecimal encoding of the input.

**Example**

Query:

```sql
select evm_hex_encode('hello')
```

Result:

| evm_hex_encode('hello') |
|:-|
| 0x68656c6c6f |