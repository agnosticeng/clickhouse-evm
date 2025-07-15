select evm_decode_event(
	[
		evm_hex_decode('0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef'),
		evm_hex_decode('0x00000000000000000000000063dfe4e34a3bfc00eb0220786238a7c6cef8ffc4'),
		evm_hex_decode('0x000000000000000000000000936c700adf05d1118d6550a3355f66e93c9476c6')
	]::Array(FixedString(32)),
	evm_hex_decode('0x0000000000000000000000000000000000000000000000000000000252e9f940'),
	['event Transfer(address indexed,address indexed,uint256)']
)