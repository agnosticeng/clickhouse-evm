### evm_decode_event

Decodes an ABI-encoded event.

**Syntax**

```sql
select evm_decode_event(topics, input_data, [dec0, dec1, ...])
```

**Parameters**

- `topics` - EVM-encoded input data of the call. [Array(FixedString(32))](https://clickhouse.com/docs/en/sql-reference/data-types/fixedstring)
- `input_data` - EVM-encoded input data of the event. [String](https://clickhouse.com/docs/en/sql-reference/data-types/string)
- `decoders` - An array of decoders for the call. A decoded can be either a [fullsig](../evm_fullsig.md) or the URL of a JSON-encoded ABI. [Array(String)](https://clickhouse.com/docs/en/sql-reference/data-types/array)

**Returned value**

- The decoded log, wrapped in a [`Result`](../error_handling.md).
  The `value` field of the [`Result`](../error_handling.md) object contains the following fields:
    - `signature` - A string representing the signature of the signature
    - `inputs` - An object containing the decoded input parameters of the function

**Example**

The below example decodes an EVM-encoded `Transfer` event from an ERC-20 contract.

Query:

```sql
select evm_decode_event(
	[
		evm_hex_decode('0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef'),
		evm_hex_decode('0x00000000000000000000000063dfe4e34a3bfc00eb0220786238a7c6cef8ffc4'),
		evm_hex_decode('0x000000000000000000000000936c700adf05d1118d6550a3355f66e93c9476c6')
	]::Array(FixedString(32)),
	evm_hex_decode('0x0000000000000000000000000000000000000000000000000000000252e9f940'),
	['event Transfer(address indexed,address indexed,uint256)']
) as res
```

Result:

| res |
|:-|
| {"value":{"signature":"Transfer(address,address,uint256)","inputs":{"arg2":"9981000000","arg0":"0x63dfe4e34a3bfc00eb0220786238a7c6cef8ffc4","arg1":"0x936c700adf05d1118d6550a3355f66e93c9476c6"}}} |

More examples are available [here](../../examples/queries/evm_decode).