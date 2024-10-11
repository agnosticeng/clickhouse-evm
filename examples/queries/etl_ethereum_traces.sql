with 
    q0 as (
        select number as n from numbers(20764111, 3)
    ),

    q1 as (
        select
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockByNumber', 
                    [evm_hex_encode_int(n), 'true'], 
                    ''
                ),
                'Tuple(
                    timestamp String,
                    transactions Array(JSON)
                )'
            ) as block,
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockReceipts', 
                    [evm_hex_encode_int(n)], 
                    ''
                ),
                'Array(JSON)'
            ) as receipts,
            JSONExtract(
                ethereum_rpc(
                    'trace_block',
                    [evm_hex_encode_int(n)], 
                    ''
                ),
                'Array(JSON)'
            ) as traces
        from q0
    ),

    q2 as (
        select 
            block,
            trace,
            block.transactions[trace.transactionPosition::UInt32+1] as tx,
            receipts[trace.transactionPosition::UInt32+1] as receipt
        from q1
        array join traces as trace
    ),

    q3 as (
        select 
            toDateTime64(evm_hex_decode_int(block.timestamp, 'Int64'), 3, 'UTC') as timestamp,
            evm_hex_decode(tx.blockHash::String) as block_hash,
            evm_hex_decode_int(tx.blockNumber::String, 'UInt64') as block_number,
            evm_hex_decode(tx.hash::String) as transaction_hash,
            evm_hex_decode_int(tx.transactionIndex::String, 'UInt32') as transaction_index,
            evm_hex_decode(tx.from::String) as transaction_from,
            evm_hex_decode_int(receipt.status::String, 'UInt8') as transaction_status,
            trace.subtraces::UInt32 as subtraces,
            trace.traceAddress::Array(UInt32) as trace_address,
            trace.type::String as type,
            trace.error::String as error,
            trace.action.callType::String as call_type,
            evm_hex_decode(trace.action.from::String) as from,
            evm_hex_decode_int(trace.action.gas::String, 'UInt64') as gas,
            evm_hex_decode(trace.action.input::String) as input,
            evm_hex_decode(trace.action.to::String) as to,  
            evm_hex_decode_int(trace.action.value::String, 'UInt256') as value,
            evm_hex_decode(trace.action.address::String) as address,  
            evm_hex_decode_int(trace.action.balance::String, 'UInt256') as balance,
            evm_hex_decode(trace.action.refundAddress::String) as refund_address,  
            evm_hex_decode(trace.action.author::String) as author,  
            trace.action.rewardType::String as reward_type,
            evm_hex_decode(trace.action.init::String) as init,
            evm_hex_decode(trace.result.address::String) as result_address,
            evm_hex_decode(trace.result.code::String) as code,
            evm_hex_decode_int(trace.result.gasUsed::String, 'UInt64') as gas_used,
            evm_hex_decode(trace.result.output::String) as output
        from q2
    )

select * from q3
