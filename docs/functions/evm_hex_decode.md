### evm_hex_decode

Decodes 0x-prefixed hexadecimal encoded strings.

**Syntax**

```sql
evm_hex_decode(str)
```

**Parameters**

- `str` - Any string. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns the decoded value of the 0x-prefixed hexadecimal input string.

**Example**

Query:

```sql
select evm_hex_decode('0x68656c6c6f')
```

Result:

| evm_hex_decode('0x68656c6c6f') |
|:-|
| hello |