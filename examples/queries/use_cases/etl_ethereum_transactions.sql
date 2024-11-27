with 
    q0 as (
        select number as n from numbers(20764111, 1)
    ),

    q1 as (
        select
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockByNumber', 
                    [evm_hex_encode_int(n), 'true'], 
                    'fail-on-error=true&fail-on-null=true'
                ),
                'value',
                'JSON'
            ) as block,
            JSONExtract(
                ethereum_rpc(
                    'eth_getBlockReceipts', 
                    [evm_hex_encode_int(n)], 
                    'fail-on-error=true&fail-on-null=true'
                ),
                'value',
                'Array(JSON)'
            ) as receipts
        from q0
    ),

    q2 as (
        select 
            toDateTime64(evm_hex_decode_int(block.timestamp::String, 'Int64'), 3, 'UTC') as timestamp,
            arrayMap(x -> tuple(
                    evm_hex_decode(JSONExtract(x, 'address', 'String')),
                    arrayMap(x -> evm_hex_decode(x), JSONExtract(x, 'storageKeys', 'Array(String)'))
                ),   
                tx.accessList::Array(String)
            ) as access_list,
            evm_hex_decode(tx.blockHash::String) as block_hash,
            evm_hex_decode_int(tx.blockNumber::String, 'UInt64') as block_number,
            evm_hex_decode_int(tx.chainId::String, 'UInt32') as chain_id,
            evm_hex_decode(tx.from::String) as from,
            evm_hex_decode_int(tx.gas::String, 'UInt64') as gas,
            evm_hex_decode_int(tx.gasPrice::String, 'UInt256') as gas_price,
            evm_hex_decode(tx.hash::String) as hash,
            evm_hex_decode(tx.input::String) as input,
            evm_hex_decode_int(tx.maxFeePerGas::String, 'UInt256') as max_fee_per_gas,
            evm_hex_decode_int(tx.maxPriorityFeePerGas::String, 'UInt256') as max_priority_fee_per_gas,
            evm_hex_decode_int(tx.nonce::String, 'UInt256') as nonce,
            evm_hex_decode(tx.r::String) as r,
            evm_hex_decode(tx.s::String) as s,
            evm_hex_decode(tx.to::String) as to,
            evm_hex_decode_int(tx.transactionIndex::String, 'UInt32') as transaction_index,
            evm_hex_decode_int(tx.type::String, 'UInt16') as type,
            evm_hex_decode(tx.v::String) as v,
            evm_hex_decode_int(tx.value::String, 'UInt256') as value,
            evm_hex_decode(tx.yParity::String) as yParity,
            evm_hex_decode(receipt.contractAddress::String) as contract_address,
            evm_hex_decode_int(receipt.cumulativeGasUsed::String, 'UInt64') as cumulative_gas_used,
            evm_hex_decode_int(receipt.effectiveGasPrice::String, 'UInt256') as effective_gas_price,
            evm_hex_decode_int(receipt.gasUsed::String, 'UInt64') as gas_used,
            evm_hex_decode(receipt.root::String) as root,
            evm_hex_decode_int(receipt.status::String, 'UInt8') as status
        from q1
        array join block.transactions[] as tx, receipts as receipt
    )

select * from q2




