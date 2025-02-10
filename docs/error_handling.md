# Error Handling

By definition, scalar ClickHouse UDFs (User-Defined Functions) can only return a single column. However, we want our function to support partial failure, meaning that an error affecting a single row should not cause the entire query to fail.

To achieve this, we return a Result type as a JSON-encoded string, encapsulating both the value and any potential error:

```json
{
    "value": ...,
    "error": "my error"
}
````

This approach ensures that errors are handled gracefully on a per-row basis while preserving overall batch execution.
