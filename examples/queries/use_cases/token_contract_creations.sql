select
    number as block_number,
    trace.transactionHash as transaction_hash,
    trace.action.from as creator,
    trace.result.address as contract_address,
    JSONExtractString(
        ethereum_rpc_call(contract_address, 'symbol()(string)', '', -1::Int64, ''),
        'value',
        'arg0'
    ) as symbol,
    JSONExtractUInt(
        ethereum_rpc_call(contract_address, 'decimals()(uint8)', '', -1::Int64, ''),
        'value',
        'arg0'
    ) as decimals
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
            type String,
            transactionHash String,
            action Tuple(
                from String
            ),
            result Tuple(
                address String
            )
        )
    )'    
) as trace
where trace.type == 'create'
having symbol != ''

