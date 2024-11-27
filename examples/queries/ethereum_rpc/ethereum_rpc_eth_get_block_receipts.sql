select
    ethereum_rpc(
        'eth_getBlockReceipts', 
        [evm_hex_encode_int(number)], 
        ''
    )
from numbers(20764111, 10)
