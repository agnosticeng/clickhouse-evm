# Fullsig

**Fullsig** is an enhanced version of the standard Ethereum function and event signatures, as defined in the [Solidity ABI specification](https://solidity-fr.readthedocs.io/fr/latest/abi-spec.html).

## Motivation

We aimed to create a **compact and self-contained** representation of function and event specifications, enabling **inline, on-the-fly decoding** of EVM logs and call traces.

While the existing definitions for functions and events are mostly sufficient, they lack some crucial information for complete self-contained decoding:

- **Function signatures** do not specify output types.
- **Event signatures** do not indicate which parameters are indexed.

## The Fullsig Definition

To address these gaps, we introduced the concept of **Fullsig** (full signature), which includes the missing details necessary for independent decoding.

Additionally, to clearly differentiate between **event** and **function** signatures, we prefix them explicitly with `event` or `function`.

## Examples

### ERC-20 `transfer` Function

```solidity
function transfer(address,uint256)(bool)
````

This includes the return type (bool), making it fully self-descriptive.

### ERC-20 `Transfer` Event

```solidity
event Transfer(address indexed,address indexed,uint256)
```

This explicitly marks indexed parameters, ensuring a complete event specification.

