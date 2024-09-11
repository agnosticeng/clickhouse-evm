with
    q0 as (
        select 
            coalesce(trace.action.input, '') as input,
            coalesce(trace.result.output, '') as output,
            'transfer(address,uint256)(bool)' as abi
        from file('./tmp/evm_blocks/*.parquet')
        array join transactions as tx 
        array join tx.traces as trace
        where left(trace.action.input, 4) = left(keccak256('transfer(address,uint256)'), 4)
    )

select * from q0 limit 10


