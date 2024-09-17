with 
    q0 as (
        select
            number,
            ethereum_rpc(
                'eth_getBlockReceipts', 
                [evm_hex_encode_uint256(number)], 
                ''
            ) as res
        from numbers(20000000, 10)
    ),

    q1 as (
        select 
            number as block_number,
            JSONExtract(
                res, 
                'Array(
                    Tuple(
                        blockHash String, 
                        transactionHash String,
                        logs Array(
                            Tuple(
                                address String,
                                topics Array(String),
                                data String,
                                logIndex UInt256
                            )
                        )
                    )
                )'
            ) as receipts
        from q0
    ),

    q2 as (
        select 
            block_number,
            receipt.blockHash as block_hash,
            receipt.transactionHash as transaction_hash,
            log.address as log_address,
            log.topics as log_topics,
            log.data as log_data,
            log.logIndex as log_index
        from q1
        array join receipts as receipt
        array join receipt.logs as log
    ),

    q3 as (
        select 
            block_number,
            transaction_hash,
            log_index,
            log_address as contract_address,
            evm_decode_event(
                arrayMap(x -> evm_hex_decode(x), log_topics),
                evm_hex_decode(log_data),
                'Transfer(indexed address, indexed address, uint256)'
            ) as evt
        from q2
        where length(log_topics) == 3
        and log_topics[1] = evm_hex_encode(keccak256('Transfer(address,address,uint256)'))
    )

select * from q3