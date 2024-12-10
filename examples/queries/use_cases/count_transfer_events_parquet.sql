select
    count(*)
from s3('https://data.agnostic.dev/ethereum-mainnet-pq/logs/*.parquet') as logs
where _file between '20240101.parquet' and '20240109.parquet'
and logs.topics[1] = keccak256('Transfer(address,address,uint256)')
