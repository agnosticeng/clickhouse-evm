with 
    q0 as (
        select number as n from numbers(20764111, 10)
    )

select
    ethereum_rpc(
        'eth_getBlockReceipts', 
        [evm_hex_encode_int(n)], 
        ''
    ),
    'Array(JSON)'
from q0
