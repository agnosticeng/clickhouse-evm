with 
    raw as (
        select 
            r.blockHash as block_hash,
            r.transactionHash as transaction_hash,
            JSONExtract(
                evm_decode_event(
                    arrayMap(x -> evm_hex_decode(x), l.topics),
                    evm_hex_decode(l.data),
                    ['Transfer(indexed address, indexed address, uint256)']
                ),
                'JSON'
            ) as evt
        from numbers(20764111, 3)  
        array join JSONExtract(
            ethereum_rpc('eth_getBlockReceipts', [evm_hex_encode_int(number)], 'fail-on-error=true&fail-on-null=true'),
            'value',
            'Array(
                Tuple(
                    blockHash String,
                    transactionHash String,
                    logs Array(
                        Tuple(
                            data String,
                            topics Array(String)
                        )
                    )
                )
            )'
        ) as r
        array join r.logs as l
        where length(l.topics) == 3
        and evm_hex_decode(l.topics[1]) = keccak256('Transfer(address,address,uint256)')
    )

select 
    block_hash,
    transaction_hash,
    evt.value.inputs.arg0::String as sender,
    evt.value.inputs.arg1::String as recipient,
    evt.value.inputs.arg2::UInt256 as amount
from raw
where evt.error is null