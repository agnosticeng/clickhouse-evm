select ethereum_rpc(
    'eth_getLogs', 
    [
        toJSONString(map(
            'address', '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640',
            'from', evm_hex_encode_uint256(20000000),
            'to', evm_hex_encode_uint256(20000005)
        ))
    ], 
    ''
)
    