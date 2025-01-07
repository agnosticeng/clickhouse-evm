## Using agnostic-clickhouse-udf with clickhouse-local

### Build bundle

```sh
make bundle
```

### Point clickhouse-local the the custom config file

```sh
clickhouse local --config ./examples/clickhouse-local-config.xml --path tmp/clickhouse
```

### Running clickhouse-server with Docker Compose

```sh
docker-compose up -d
```

### Installing a bundle in a running clickhouse-server container

```sh
su - clickhouse -c "wget -qO- https://github.com/agnosticeng/agnostic-clickhouse-udf/releases/download/v0.0.3/agnostic-clickhouse-udf_0.0.3_linux_amd64_v3.tar.gz | tar xvz -C /"
for f in /var/lib/clickhouse/user_defined/*.sql; do clickhouse client --queries-file $f; done
```

Here are the SQL function signatures with proper ClickHouse types:

1. **ethereum_rpc_call**
```sql
CREATE FUNCTION ethereum_rpc_call AS (
    to String,              -- Contract address (0x...)
    fullsig String,         -- Function signature (e.g. "transfer(address,uint256)")
    data Array(String),     -- Call parameters as JSON array
    block_number Int64,     -- Block number (-2 for latest, -3 for finalized etc)
    endpoint String         -- RPC endpoint URL
) -> Nullable(String);      -- Returns JSON response
```

2. **ethereum_rpc**
```sql
CREATE FUNCTION ethereum_rpc AS (
    method String,          -- RPC method name
    params String,          -- RPC parameters as JSON
    endpoint String         -- RPC endpoint URL
) -> Nullable(String);      -- Returns JSON response
```

3. **evm_decode_call**
```sql
CREATE FUNCTION evm_decode_call AS (
    input String,           -- Call input data (0x...)
    output String,          -- Call output data (0x...)
    abi Array(String)       -- Array of ABI provider strings
) -> Nullable(String);      -- Returns decoded data as JSON
```

4. **evm_decode_event**
```sql
CREATE FUNCTION evm_decode_event AS (
    topics Array(String),   -- Event topics array (0x...)
    data String,           -- Event data (0x...)
    abi Array(String)      -- Array of ABI provider strings
) -> Nullable(String);     -- Returns decoded event as JSON
```

5. **evm_descriptor_from_fullsig**
```sql
CREATE FUNCTION evm_descriptor_from_fullsig AS (
    fullsig String,        -- Full function/event signature
    type String           -- Type ('event' or 'function')
) -> Nullable(String);    -- Returns ABI descriptor as JSON
```

6. **evm_signature_from_descriptor**
```sql
CREATE FUNCTION evm_signature_from_descriptor AS (
    event_descriptor String  -- ABI descriptor as JSON
) -> Nullable(String);      -- Returns signature info as JSON {selector, signature, fullsig}
```

7. **keccak256**
```sql
CREATE FUNCTION keccak256 AS (
    str String             -- Input string/bytes
) -> FixedString(32);     -- Returns 32-byte Keccak256 hash
```

8. **ethereum_rpc_filter** (Table Function)
```sql
CREATE FUNCTION ethereum_rpc_filter AS (
    filter String          -- Filter configuration as JSON
) -> Table (
    result String         -- Filter results as JSON
);

-- Usage example:
SELECT *
FROM ethereum_rpc_filter('{"fromBlock": "latest", "address": ["0x..."]}')
SETTINGS poll_method='eth_getFilterChanges', poll_interval=1;
```

Common Usage Examples:
```sql
-- Call contract function
SELECT ethereum_rpc_call(
    '0x1234...', 
    'balanceOf(address)', 
    ['0x5678...'],
    -2, 
    'https://eth-mainnet.g.alchemy.com/v2/key'
);

-- Decode event
SELECT evm_decode_event(
    ['0x...', '0x...'], 
    '0x...', 
    ['https://api.etherscan.io/api']
);

-- Calculate Keccak256
SELECT hex(keccak256('Hello'));
```

Notes:
- All hex strings should be prefixed with '0x'
- Block numbers: -2 = latest, -3 = finalized, -1 = pending, 0 = earliest
- JSON responses typically include {value: result} or {error: message} structure
- ABI providers can be URLs or inline ABI JSON strings
- All functions support batch processing for better performance
