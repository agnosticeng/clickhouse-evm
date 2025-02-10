### ethereum_rpc

Call [Ethereum RPC methods](https://ethereum.org/fr/developers/docs/apis/json-rpc/).

**Syntax**

```sql
ethereum_rpc(method, [param0, param1, ...], endpoint)
```

**Parameters**

- `method` - Any RPC method supported by the RPC endpoint. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `params` - An array of JSON-encoded parameters for the RPC method. [Array](https://clickhouse.com/docs/en/sql-reference/data-types/array)
- `endpoint` - An RPC endpoint. Can be left blank to use default endpoint. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- The response to the RPC call, wrapped in a [`Result`](../error_handling.md).

**Example**

The below example fetch the transaction receipts for a range of 10 blocks from the RPC node.
More examples can be found [here](examples/queries/ethereum_rpc).

Query:

```sql
select
    ethereum_rpc(
        'eth_getBlockTransactionCountByNumber', 
        [evm_hex_encode_int(number)], 
        'https://eth.llamarpc.com'
    ) as tx_count
from numbers(20764111, 10)
```

Result:

| tx_count |
|:-|
| {"value":"0x90"} |
| {"value":"0xfe"} |
| {"value":"0x9c"} |
| {"value":"0x84"} |
| {"value":"0x9b"} |
| {"value":"0x89"} |
| {"value":"0xa9"} |
| {"value":"0xb1"} |
| {"value":"0xc4"} |
| {"value":"0x94"} |