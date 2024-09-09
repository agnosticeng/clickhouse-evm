select 
    count(*)
from file('./tmp/evm_blocks/*.parquet')