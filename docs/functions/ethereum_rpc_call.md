### ethereum_rpc_call

Call a contract [view function](https://docs.soliditylang.org/en/latest/contracts.html#view-functions) using the [eth_call](https://docs.alchemy.com/reference/eth-call) RPC method.

**Syntax**

```sql
ethereum_rpc_call(contract_address, fullsig, data, block_number, endpoint)
```

**Parameters**

- `contract_address` - The address of the contract on which to call the view function. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `fullsig` - The [fullsig](../evm_fullsig.md) of the function to call. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `data` - The input data of the function passed as a JSON object or Array. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `block_number` - The block number at which the state of the execution engine must be set before the function is called. [Int64](https://clickhouse.com/docs/en/sql-reference/data-types/int-uint)
- `endpoint` - An RPC endpoint. Can be left blank to use default endpoint. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- The response to the RPC call, wrapped in a [`Result`](../error_handling.md).

**Example**

The below example fetch balance of [USDT](https://etherscan.io/token/0xdac17f958d2ee523a2206206994597c13d831ec7) tokens of the [Kraken 4 wallet](https://etherscan.io/address/0x267be1c1d684f78cb4f6a176c4911b741e4ffdc0) at latest block.

Query:

```sql
select 
    ethereum_rpc_call(
        '0xdac17f958d2ee523a2206206994597c13d831ec7', 
        'function balanceOf(address)(uint256)',
        toJSONString(['0x267be1c1d684f78cb4f6a176c4911b741e4ffdc0']),
        -1::Int64,
        'https://eth.llamarpc.com#fail-on-retryable-error=true&fail-on-null=true'
)
```

Result:

| balance |
|:-|
| {"value":{"arg0":"2097761566"}} |



