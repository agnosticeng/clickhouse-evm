with 
    block_numbers as (
        select 
            generate_series as n 
        from generate_series(
            20000000,
            20000010
        )
    ),

    q0 as (
        select
            n as block_number,
            JSONExtract(
                ethereum_rpc('eth_getBlockByNumber', [evm_hex_encode_int(n), 'false'], ''),
                'Tuple(
                    timestamp String
                )'
            ) as block,
            JSONExtract(
                ethereum_rpc('eth_getBlockReceipts', [evm_hex_encode_int(n)], ''),
                'Array(
                    Tuple(
                        blockHash String,
                        blockNumber String,
                        from String,
                        status String,
                        transactionHash String,
                        transactionIndex String
                    )
                )'
            ) as receipts,
            JSONExtract(
                ethereum_rpc('trace_block', [evm_hex_encode_int(n)], ''),
                'Array(
                    Tuple(
                        transactionPosition UInt32,
                        subtraces UInt32,
                        traceAddress Array(UInt32),
                        type String,
                        error String,
                        callType String,
                        action Tuple(
                            from String,
                            gas String,
                            input String,
                            to String,
                            value String,
                            address String,
                            balance String,
                            refundAddress String,
                            author String,
                            rewardType String,
                            init String,
                        ),
                        result Tuple(
                            address String,
                            code String,
                            gasUsed String,
                            output String
                        )
                    )
                )'
            ) as traces
        from block_numbers
    ),

    q1 as (
        select
            arrayMap(
                x -> tuple(
                    toDateTime64(evm_hex_decode_int(block.timestamp, 'Int64'), 3, 'UTC') as timestamp, 
                    evm_hex_decode(receipts[x.transactionPosition].blockHash) as block_hash,
                    evm_hex_decode_int(receipts[x.transactionPosition].blockNumber::String, 'UInt64') as block_number,
                    evm_hex_decode(receipts[x.transactionPosition].from::String) as transaction_from,
                    evm_hex_decode_int(receipts[x.transactionPosition].status::String, 'UInt8') as transaction_status,
                    evm_hex_decode(receipts[x.transactionPosition].transactionHash::String) as transaction_hash,
                    evm_hex_decode_int(receipts[x.transactionPosition].transactionIndex::String, 'UInt32') as transaction_index,
                    x.subtraces as subtraces,
                    x.traceAddress as trace_address,
                    x.type as type,
                    x.error as error,
                    x.callType as call_type,
                    evm_hex_decode(x.action.from::String) as from,
                    evm_hex_decode_int(x.action.gas::String, 'UInt64') as gas,
                    evm_hex_decode(x.action.input::String) as input,
                    evm_hex_decode(x.action.to::String) as to,  
                    evm_hex_decode_int(x.action.value::String, 'UInt256') as value,
                    evm_hex_decode(x.action.address::String) as address,  
                    evm_hex_decode_int(x.action.balance::String, 'UInt256') as balance,
                    evm_hex_decode(x.action.refundAddress::String) as refund_address,  
                    evm_hex_decode(x.action.author::String) as author,  
                    x.action.rewardType::String as reward_type,
                    evm_hex_decode(x.action.init::String) as init,
                    evm_hex_decode(x.result.address::String) as result_address,
                    evm_hex_decode(x.result.code::String) as result_code,
                    evm_hex_decode_int(x.result.gasUsed::String, 'UInt64') as gas_used,
                    evm_hex_decode(x.result.output::String) as output
                ), 
                traces
            ) as traces
        from q0
    )

select 
    untuple(t)
from q1
array join traces as t







