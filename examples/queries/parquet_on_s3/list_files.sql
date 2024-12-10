select 
    _file 
from s3('https://data.agnostic.dev/ethereum-mainnet-pq/blocks/*.parquet', 'One') 
settings remote_filesystem_read_prefetch=false