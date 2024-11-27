select
    ethereum_rpc(
        'eth_getBlockReceipts', 
        [evm_hex_encode_int(number)], 
        'https://eth.llamarpc.com#fail-on-null=true'
    )
from numbers(50000000, 1)
