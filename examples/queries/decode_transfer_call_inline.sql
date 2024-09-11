with 
    q0 as (
        select 
            number as block_number,
            tx.hash as tx_hash,
            tx.from as from,
            trace.trace_address as trace_address,
            evm_decode_call(
                coalesce(trace.action.input, ''),
                coalesce(trace.result.output, ''),
                'transfer(address,uint256)(bool)'
            ) as call
        from file('./tmp/evm_blocks/*.parquet')
        array join transactions as tx 
        array join tx.traces as trace
        where left(trace.action.input, 4) = left(keccak256('transfer(address,uint256)'), 4)
    ),

    q1 as (
        select 
            block_number,
            evm_hex_encode(tx_hash) as tx_hash,
            trace_address as trace_address,
            evm_hex_encode(from) as sender,
            evm_hex_encode(JSON_VALUE(call, '$.inputs.arg0')) as recipient,
            JSON_VALUE(call, '$.inputs.arg1') as amount
        from q0
    )

select * from q1 limit 10;

