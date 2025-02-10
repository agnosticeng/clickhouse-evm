### keccak256

Compute the [keccak256](https://github.com/ethereum/eth-hash) hash of a String.

**Syntax**

```sql
keccak256(str)
```

**Parameters**

- `str` - Any string. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- Returns the keccack256 hash of the input str

**Example**

This example computes the [keccak256](https://github.com/ethereum/eth-hash) hash of the string `Transfer(address,address,uint256)` and encodes it as 0x-prefixed hex.

Query:

```sql
select evm_hex_encode(keccak256('Transfer(address,address,uint256)'))
```

Result:

```response
"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
```

