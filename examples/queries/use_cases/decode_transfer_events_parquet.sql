with 
    raw as (
        select
            evm_hex_encode(logs.block_hash) as block_hash,
            evm_hex_encode(logs.transaction_hash) as transaction_hash,
            JSONExtract(
                evm_decode_event(
                    logs.topics::Array(String),
                    logs.data::String,
                    ['event Transfer(address indexed, address indexed, uint256)']
                ),
                'JSON'
            ) as evt
        from s3('https://data.agnostic.dev/ethereum-mainnet-pq/logs/*.parquet') as logs
        where _file between '20240101.parquet' and '20240109.parquet'
        and logs.topics[1] = keccak256('Transfer(address,address,uint256)')
        limit 100
    )

select 
    block_hash,
    transaction_hash,
    evt.value.inputs.arg0::String as sender,
    evt.value.inputs.arg1::String as recipient,
    evt.value.inputs.arg2::UInt256 as amount
from raw
where evt.error is null
