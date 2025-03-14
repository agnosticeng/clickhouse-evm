### ethereum_decode_tx

Decodes an RLP-encoded [Ethereum transaction](https://ethereum.org/en/developers/docs/transactions/).

**Syntax**

```sql
select evm_decode_tx(input_data)
```

**Parameters**

- `input_data` - RLP-encoded [Ethereum transaction](https://ethereum.org/en/developers/docs/transactions/). [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)

**Returned value**

- The decoded [Ethereum transaction](https://ethereum.org/en/developers/docs/transactions/) object, wrapped in a [`Result`](../error_handling.md).

**Example**

Query:

```sql
select ethereum_decode_tx(evm_hex_decode('0x02f8b1018201f08305c1d58439c652a98301482094dac17f958d2ee523a2206206994597c13d831ec780b844a9059cbb0000000000000000000000005ff90de9d2aedb02c2924011dfbb6d12d35c81180000000000000000000000000000000000000000000000000000000003e96f30c001a0e939a9cb318e770a96b8dad62cde0be0eb1ed9c9ee1331690d0f343cc2e7ca8da01172360c496760cc7fee1e538951728db6ded66aa8cf6bf710d4f968c78ac1e6
')) as tx
```

Result:

| tx |
|:-|
| {"value":{"type":"0x2","chainId":"0x1","nonce":"0x1f0","to":"0xdac17f958d2ee523a2206206994597c13d831ec7","gas":"0x14820","gasPrice":null,"maxPriorityFeePerGas":"0x5c1d5","maxFeePerGas":"0x39c652a9","value":"0x0","input":"0xa9059cbb0000000000000000000000005ff90de9d2aedb02c2924011dfbb6d12d35c81180000000000000000000000000000000000000000000000000000000003e96f30","accessList":[],"v":"0x1","r":"0xe939a9cb318e770a96b8dad62cde0be0eb1ed9c9ee1331690d0f343cc2e7ca8d","s":"0x1172360c496760cc7fee1e538951728db6ded66aa8cf6bf710d4f968c78ac1e6","yParity":"0x1","hash":"0xb043de7481fc43efae8933482afa79b8c6a990be8c03b241d51127aa3ebd7fa1"}} |