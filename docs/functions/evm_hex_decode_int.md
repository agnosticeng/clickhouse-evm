### evm_hex_decode_int

Decodes 0x-prefixed hexadecimal encoded strings to integers.

**Syntax**

```sql
evm_hex_decode(str)
```

**Parameters**

- `str` - Any string. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `type` - Any [integer type](https://clickhouse.com/docs/en/sql-reference/data-types/int-uint) name. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns the decoded value of input 0x-prefixed hexadecimal input integer.

**Example**

Query:

```sql
select evm_hex_decode_int('0x7b', 'Int64')
```

Result:

| evm_hex_decode_int('0x7b', 'Int64') |
|-:|
| 123 |