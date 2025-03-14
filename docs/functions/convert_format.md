### convert_format

Convert from one of the supported text serialization format to another.

Supported formats are:
- **JSON**
- **YAML**
- **TOML**

**Syntax**

```sql
select convert_format(from_format, to_format, input_data)
```

**Parameters**

- `from_format` - The format of the input data. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `to_format` - The expected output format. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `input_data` - A string containing a valid serialized object in the input format. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- The original object serialized as a valid output format value

**Example**

Query:

```sql
select convert_format('YAML', 'JSON', $heredoc$
    root:
        key1: val
        arr: 
            - item1
            - item2
$heredoc$) as v
```

Result:

| v |
|:-|
| {"root":{"arr":["item1","item2"],"key1":"val"}} |