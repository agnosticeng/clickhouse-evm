with 
    q0 as (
        select 
            block_number,
            evm_hex_encode(transaction_hash) as tx_hash,
            log_index,
            address,
            evm_decode_event(
                topics::Array(String),
                data::String,
                'Transfer(indexed address, indexed address, uint256)'
            ) as evt
        from s3('https://data.agnostic.dev/ethereum-mainnet-pq/logs/*.parquet', 'Parquet')
        where _file >= '20241001'
        and address = evm_hex_decode('0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce')
        and length(topics) == 3
        and topics[1] = keccak256('Transfer(address,address,uint256)')
    )

select 
    block_number,
    tx_hash,
    log_index,
    evm_hex_encode(address) as token,
    evm_hex_encode(JSON_VALUE(evt, '$.inputs.arg0')) as sender,
    evm_hex_encode(JSON_VALUE(evt, '$.inputs.arg1')) as recipient,
    JSON_VALUE(evt, '$.inputs.arg2') as amount    
from q0

