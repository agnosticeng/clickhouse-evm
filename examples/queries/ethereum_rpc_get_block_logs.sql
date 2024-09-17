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
                                data String
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
            log.data as log_data
        from q1
        array join receipts as receipt
        array join receipt.logs as log
    )

select * from q2