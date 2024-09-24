with 
    q0 as (
        select
            number,
            ethereum_rpc(
                'trace_block', 
                [evm_hex_encode_uint256(number)], 
                ''
            ) as res
        from numbers(6082465, 10)
    ),

    q1 as (
        select 
            number as block_number,
            JSONExtract(
                res, 
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
            ) as traces
        from q0
    ),

    q2 as (
        select
            block_number,
            trace.transactionHash as transaction_hash,
            trace.action.from as creator,
            trace.result.address as contract_address
        from q1
        array join traces as trace
        where trace.type == 'create'
    ),

    q3 as (
        select
            q2.*,
            JSONExtractString(
                ethereum_rpc_call(q2.contract_address, 'symbol()(string)', '', -1::Int64, ''),
                'data',
                'arg0'
            ) as symbol,
            JSONExtractUInt(
                ethereum_rpc_call(q2.contract_address, 'decimals()(uint8)', '', -1::Int64, ''),
                'data',
                'arg0'
            ) as decimals
        from q2
        having symbol != ''
    )

select * from q3