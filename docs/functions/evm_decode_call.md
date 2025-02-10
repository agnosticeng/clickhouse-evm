### evm_decode_call

Decodes an ABI-encoded function call.

**Syntax**

```sql
select evm_decode_call(input_data, output_data, [dec0, dec1, ...])
```

**Parameters**

- `input_data` - EVM-encoded input data of the call. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `output_data` - EVM-encoded output data of the call. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `decoders` - An array of decoders for the call. A decoded can be either a [fullsig](../evm_fullsig.md) or the URL of a JSON-encoded ABI. [Array(String)](https://clickhouse.com/docs/en/sql-reference/data-types/array)

**Returned value**

- The decoded function call, wrapped in a [`Result`](../error_handling.md).
  The `value` field of the [`Result`](../error_handling.md) object contains the following fields:
    - `signature` - A string representing the signature of the signature
    - `inputs` - An object containing the decoded input parameters of the function
    - `outputs` - An object containg the output parameters of the function

**Example**

The below example decodes an EVM-encoded `transfer` call trace from an ERC-20 contract.

Query:

```sql
select evm_decode_call(
	evm_hex_decode('0xa9059cbb0000000000000000000000005e6cb68740e2ade791a083d19d339580d510948000000000000000000000000000000000000000000000000000000000065d7c70'),
	evm_hex_decode('0x0000000000000000000000000000000000000000000000000000000000000001'),
	['function transfer(address,uint256)(bool)']
) as res
```

Result:

| res |
|:-|
| {"value":{"signature":"transfer(address,uint256)","inputs":{"arg0":"0x5e6cb68740e2ade791a083d19d339580d5109480","arg1":"106790000"},"outputs":{"arg0":true}}} |

More examples are available [here](../../examples/queries/evm_decode).