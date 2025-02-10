# JSON-RPC Client

All functions that make RPC calls to Ethereum nodes utilize the same in-house JSON-RPC client.

By default, this client automatically batches calls to nodes, as described in the Geth documentation. Since ClickHouse processes UDFs (User-Defined Functions) with full blocks of data, batching is naturally aligned with ClickHouse’s execution model, improving efficiency and reducing RPC overhead.

## Passing options

The client offers various configurable options, which can be set in two ways:

- Globally for all UDF instances using environment variables or command-line flags.
- Per-call basis using URL hash parameters, e.g.:

```
https://eth.llamarpc.com#max-batch-size=50

```

___

### Configuration Options



| URL Hash Parameter         | Environment Variable              | Command-Line Flag         | Type   | Default Value | Description |
|---------------------------|----------------------------------|---------------------------|--------|---------------|-------------|
|                             | `ETHEREUM_RPC_ENDPOINT`         | `endpoint`                 | string | *(required)*   | The JSON-RPC endpoint for sending calls. |
| `max-batch-size`           | `ETHEREUM_RPC_MAX_BATCH_SIZE`   | `max-batch-size`          | int64  | `200`         | Maximum number of calls per batch. |
| `max-concurrent-requests`  | `ETHEREUM_RPC_MAX_CONCURRENT_REQUESTS` | `max-concurrent-requests` | int64  | `5`           | Maximum number of concurrent outgoing RPC calls. |
| `disable-batch`            | `ETHEREUM_RPC_DISABLE_BATCH`    | `disable-batch`           | bool   | `false`       | Disables batching, sending one RPC request per row instead. |
| `fail-on-error`            | `ETHEREUM_RPC_FAIL_ON_ERROR`    | `fail-on-error`           | bool   | `false`       | Fails the entire batch if at least one RPC call encounters an error. |
| `fail-on-retryable-error`  | `ETHEREUM_RPC_FAIL_ON_RETRYABLE_ERROR` | `fail-on-retryable-error` | bool   | `false`       | Similar to `fail-on-error`, but only fails on **retryable** errors (which vary by blockchain). For example, **Arbitrum** nodes may temporarily return `intrinsic gas too low` under certain conditions. |
| `fail-on-null`             | `ETHEREUM_RPC_FAIL_ON_NULL`     | `fail-on-null`            | bool   | `false`       | Fails the batch if any RPC call returns a `null` response. |

___

This JSON-RPC client is designed to optimize performance while providing flexibility in handling errors and batching behavior. 🚀
