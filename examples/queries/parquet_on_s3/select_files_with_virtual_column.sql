select 
    count(*) 
from s3('https://data.agnostic.dev/ethereum-mainnet-pq/blocks/*.parquet') 
where _file between '20240101.parquet' and '20240109.parquet'