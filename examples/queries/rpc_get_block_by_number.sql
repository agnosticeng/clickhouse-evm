with 
    q0 as (
        select
            number,
            ethereum_rpc(
                'eth_getBlockByNumber', 
                [evm_hex_encode_uint256(number), 'false'], 
                ''
            ) as res
        from numbers(20000000, 10)
    )

select 
    number,
    JSON_VALUE(res, '$.hash') as hash
from q0