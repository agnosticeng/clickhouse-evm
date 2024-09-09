with 
    q0 as (
        select 
            number as block_number,
            tx.hash as tx_hash,
            log.log_index as log_index,
            log.address as contract_address,
            evm_decode_event(
                arrayMap(x -> coalesce(x, ''), log.topics),
                coalesce(log.data, '')
            ) as evt
        from file('./tmp/evm_blocks/*.parquet')
        array join transactions as tx 
        array join tx.receipt.logs as log
        where length(log.topics) == 3
        and log.topics[1] = keccak256('Transfer(address,address,uint256)')
    ),

    q1 as (
        select 
            block_number,
            evm_hex_encode(tx_hash) as tx_hash,
            log_index,
            evm_hex_encode(contract_address) as token,
            evm_hex_encode(JSON_VALUE(evt, '$.inputs.from')) as sender,
            evm_hex_encode(JSON_VALUE(evt, '$.inputs.to')) as recipient,
            JSON_VALUE(evt, '$.inputs.value') as amount
        from q0
    )

select * from q1 limit 10;
