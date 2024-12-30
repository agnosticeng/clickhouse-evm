select 
    JSONExtractString(
        ethereum_rpc_call(
            '0xdAC17F958D2ee523a2206206994597C13D831ec7', 
            'function symbol()(string)', 
            '', 
            -1::Int64, 
            'https://eth.llamarpc.com'
        ),
        'value',
        'arg0'
    ) as symbol,
    JSONExtractUInt(
        ethereum_rpc_call(
            '0xdAC17F958D2ee523a2206206994597C13D831ec7', 
            'function decimals()(uint8)', 
            '', 
            -1::Int64, 
            'https://eth.llamarpc.com'
        ),
        'value',
        'arg0'
    ) as decimals
