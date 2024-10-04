with 
    q0 as (
        select number as n from numbers(20764111, 1)
    ),

    q1 as (
        select
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockByNumber', 
                    [evm_hex_encode_int(n), 'false'], 
                    ''
                ),
                'JSON'
            ) as block,
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockReceipts', 
                    [evm_hex_encode_int(n)], 
                    ''
                ),
                'Array(JSON)'
            ) as receipts
        from q0
    ),

    q2 as (
        select
            toDateTime64(evm_hex_decode_int(block.timestamp::String, 'Int64'), 3, 'UTC') as timestamp,
            evm_hex_decode(receipt.blockHash::String) as block_hash,
            evm_hex_decode_int(receipt.blockNumber::String, 'UInt64') as block_number,
            evm_hex_decode(receipt.from::String) as transaction_from,
            evm_hex_decode_int(receipt.status::String, 'UInt8') as transaction_status,
            evm_hex_decode(receipt.transactionHash::String) as transaction_hash,
            evm_hex_decode_int(receipt.transactionIndex::String, 'UInt32') as transaction_index,
            toBool(log.removed::String) as removed,
            evm_hex_decode_int(log.logIndex::String, 'UInt32') as log_index,
            evm_hex_decode(log.address::String) as address,
            evm_hex_decode(log.data::String) as data,
            arrayMap(x -> evm_hex_decode(x), log.topics::Array(String)) as topics
        from q1
        array join receipts as receipt
        array join receipt.logs[] as log
    )

select * from q2












