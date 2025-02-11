select
    ethereum_rpc(
        'eth_getBlockReceipts', 
        [evm_hex_encode_int(number)], 
        'https://eth.llamarpc.com'
    )
from numbers(20764111, 10)
