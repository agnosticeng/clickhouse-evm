### evm_signature_from_descriptor

Compute the selector, signature and [fullsig](../evm_fullsig.md) from a [JSON-encoded ABI field descriptor](https://docs.soliditylang.org/en/latest/abi-spec.html#json).

**Syntax**

```sql
evm_signature_from_descriptor(json_descptor)
```

**Parameters**

- `str` - A JSON-encoded ABI field descriptor. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns a string containing the JSON encoding of an object with the below fields:
    - selector
    - signature
    - fullsig
    

**Example**

Query:

```sql
select evm_signature_from_descriptor($str$
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "from",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "to",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Transfer",
        "type": "event"
    }
$str$) as sig
```

Result:

| sig |
|:-|
| {"value":{"selector":"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","signature":"Transfer(address,address,uint256)","fullsig":"event Transfer(address indexed,address indexed,uint256)"}} |


