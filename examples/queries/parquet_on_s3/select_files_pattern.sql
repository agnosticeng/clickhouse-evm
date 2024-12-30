select 
    count(*) 
from s3('https://data.agnostic.dev/ethereum-mainnet-pq/blocks/2024010*.parquet')