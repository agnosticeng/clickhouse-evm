with 
    (
        select JSONExtract(
            ethereum_rpc(
                'eth_getLogs', 
                [
                    toJSONString(map(
                        'address', '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640',
                        'fromBlock', evm_hex_encode_int(20000000),
                        'toBlock', evm_hex_encode_int(20000010)
                    ))
                ], 
                ''
            ),
            'value',
            'Array(
                Tuple(
                    address String,
                    topics Array(String),
                    data String,
                    blockNumber String,
                    blockHash String,
                    transactionHash String,
                    transactionIndex String,
                    logIndex String,
                    removed Bool
                )
            )'
            )
    ) as raw_logs,

    logs as (
        select 
            evm_hex_decode_int(log.blockNumber, 'UInt64') as block_number,
            log.blockHash as block_hash,
            log.transactionHash as transaction_hash,
            evm_hex_decode_int(log.transactionIndex, 'UInt16') as transaction_index,
            evm_hex_decode_int(log.logIndex, 'UInt16') as log_index,
            log.address as address,
            arrayMap(t -> t, log.topics) as topics,
            log.data as data,
            toBool(log.removed) as removed
        from system.one
        array join raw_logs as log
    )

select * from logs

