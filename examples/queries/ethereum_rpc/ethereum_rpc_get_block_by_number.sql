select
    number,
    ethereum_rpc(
        'eth_getBlockByNumber', 
        [evm_hex_encode_int(number), 'false'], 
        ''
    ) as res
from numbers(20000000, 10)
