# clickhouse-evm

**EVM-focused UDFs for ClickHouse – Accelerate Ethereum & EVM-compatible blockchain analytics**

## Overview

`clickhouse-evm` is a collection of high-performance **user-defined functions (UDFs)** that extend [ClickHouse](https://clickhouse.com/) with capabilities tailored for **Ethereum Virtual Machine (EVM)** data processing.

Whether you're building blockchain explorers, indexing on-chain data, or running deep analytics on Ethereum or other EVM-compatible chains, this project brings native decoding, parsing, and querying support into your ClickHouse workflows.

## ✨ Features

- 🧠 Decode EVM calldata and logs directly within ClickHouse
- 🔐 ABI decoding for function inputs, event topics, and return values
- 🔄 Keccak-256 hashing UDF for topic and selector lookups
- 🧱 Utility functions for working with Ethereum data types (e.g., addresses, uint256, bytes, etc.)
- ⚡ **Fast, optimized RPC calls** to EVM-compatible nodes directly from ClickHouse queries
- 🚀 Speeds up on-chain data analysis by reducing external parsing overhead

## 📦 Use Cases

- Quickly extract function parameters from calldata in a ClickHouse query
- Filter and parse smart contract events by ABI signature
- Analyze token transfers, contract interactions, or DeFi protocol data at scale
- **Query live on-chain data** through efficient RPC calls from within ClickHouse
- Integrate with your existing ClickHouse-based blockchain indexing pipeline

## 📦 Artifact: The Bundle

The output of the build process is distributed as a **compressed archive** called a **bundle**. This bundle includes everything needed to deploy and use the UDFs in ClickHouse.

### 📁 Bundle Contents

Each bundle contains:

- 🧩 **Standalone binary** implementing the native UDFs (compiled with ClickHouse compatibility)
- ⚙️ **ClickHouse configuration files** (`.xml`) to register each native UDF
- 📝 **SQL files** for SQL-based UDFs (used for lightweight functions where SQL outperforms compiled code)

### 📦 Bundle Usage

#### 🛠️ Build the Bundle

```sh
make bundle              # Build for native execution
GOOS=linux make bundle   # Cross-compile for use in Docker (Linux target)
```

This will:

- Generate the bundle directory at `tmp/bundle/`
- Create a compressed archive at `tmp/bundle.tar.gz`

The internal file structure of the bundle reflects the default layout of a basic ClickHouse installation.  
As a result, **decompressing the archive at the root of a ClickHouse server filesystem should "just work"** with no additional path configuration.

---

#### ▶️ Run with `clickhouse-local`

```sh
clickhouse local --config ./examples/clickhouse-local-config.xml --path tmp/clickhouse
```

This runs ClickHouse in local mode using the provided config and a temporary storage path.

---

#### 🐳 Run with `clickhouse-server` in Docker

```sh
docker compose up -d
```

This launches a ClickHouse server inside a Docker container using the configuration and UDFs from the bundle.
