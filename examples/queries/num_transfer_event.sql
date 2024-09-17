select 
    count(*)
from file('./tmp/evm_blocks/*.parquet')
array join transactions as tx 
array join tx.receipt.logs as log
where length(log.topics) == 3
and log.topics[1] = keccak256('Transfer(address,address,uint256)')


