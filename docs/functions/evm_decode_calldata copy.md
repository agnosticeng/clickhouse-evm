### ethereum_decode_tx

Decodes raw RLP-encoded Ethereum transactions.

**Syntax**

```sql
select ethereum_decode_tx(tx)
```

**Parameters**

- `tx` - RLP-encoded Ethereum transaction. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- The decoded transaction, wrapped in a [`Result`](../error_handling.md).
  The `value` field of the [`Result`](../error_handling.md) object contains the decoded transaction as described [here](https://ethereum.org/en/developers/docs/transactions/).

**Example**

Query:

```sql
select ethereum_decode_tx(evm_hex_decode('0x02f870018314723580842fab4977825208944675c7e5baafbffbca748158becba61ef3b0a26387b12867729c43ef80c001a09fd26b54c6e097b0f71595bfa9809dec1ada35c73a164d52c44440cdf1c2811ca04136f99e25386f915aaf5ed8b88f19ead42b177131bf562c236402426a064ba2')) as tx
```

Result:

| tx |
|:-|
| {"value":{"type":"0x2","chainId":"0x1","nonce":"0x147235","to":"0x4675c7e5baafbffbca748158becba61ef3b0a263","gas":"0x5208","gasPrice":null,"maxPriorityFeePerGas":"0x0","maxFeePerGas":"0x2fab4977","value":"0xb12867729c43ef","input":"0x","accessList":[],"v":"0x1","r":"0x9fd26b54c6e097b0f71595bfa9809dec1ada35c73a164d52c44440cdf1c2811c","s":"0x4136f99e25386f915aaf5ed8b88f19ead42b177131bf562c236402426a064ba2","yParity":"0x1","hash":"0xdf28196d5e1193da59731cdcab114a3bd051da9be117f164db877ae542885893"}} |