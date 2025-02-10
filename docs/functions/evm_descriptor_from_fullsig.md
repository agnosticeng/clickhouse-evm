### evm_decriptor_from_fullsig

Tries to parse the passed [fullsig](../evm_fullsig.md) and return a [JSON-encoded ABI field descriptor](https://docs.soliditylang.org/en/latest/abi-spec.html#json).

**Syntax**

```sql
evm_descriptor_from_fullsig(fullsig)
```

**Parameters**

- `str` - An ABI field [fullsig](../evm_fullsig.md). [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns a string containing a [JSON-encoded ABI field descriptor](https://docs.soliditylang.org/en/latest/abi-spec.html#json) correspinding to the passed fullsig.

**Example**

Query:

```sql
select evm_descriptor_from_fullsig('event Transfer(address indexed, address indexed, uint256)')
```

Result:

| desc |
|:-|
| {"value":{"type":"event","name":"Transfer","inputs":[{"name":"arg0","type":"address","internalType":"address","indexed":true},{"name":"arg1","type":"address","internalType":"address","indexed":true},{"name":"arg2","type":"uint256","internalType":"uint256"}]}} |