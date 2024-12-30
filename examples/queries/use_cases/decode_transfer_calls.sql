with 
    raw as (
        select
            trace.blockHash as block_hash,
            trace.transactionHash as transaction_hash,
            trace.action.from as from,
            JSONExtract(
                evm_decode_call(
                    evm_hex_decode(trace.action.input),
                    evm_hex_decode(trace.result.output),
                    ['function transfer(address,uint256)(bool)']
                ),
                'JSON'
            ) as call
        from numbers(6082465, 10)
        array join JSONExtract(
            ethereum_rpc(
                'trace_block', 
                [evm_hex_encode_int(number)], 
                'fail-on-error=true&fail-on-null=true'
            ),
            'value',
            'Array(
                Tuple(
                    blockHash String,
                    transactionHash String,
                    type String,
                    action Tuple(
                        from String,
                        input String
                    ),
                    result Tuple(
                        output String
                    )
                )
            )'    
        ) as trace
        where left(evm_hex_decode(trace.action.input), 4) = left(keccak256('transfer(address,uint256)'), 4)
    )

select
    block_hash,
    transaction_hash,
    from as sender,
    call.value.inputs.arg0::String as recipient,
    call.value.inputs.arg1::UInt256 as amount
from raw
where call.error is null