select 
    count(*)
from numbers(20764111, 3)  
array join JSONExtract(
    ethereum_rpc('eth_getBlockReceipts', [evm_hex_encode_int(number)], 'fail-on-error=true&fail-on-null=true'),
    'value',
    'Array(
        Tuple(
            logs Array(
                Tuple(
                    topics Array(String)
                )
            )
        )
    )'
) as r
array join r.logs as l
where length(l.topics) == 3
and evm_hex_decode(l.topics[1]) = keccak256('Transfer(address,address,uint256)')