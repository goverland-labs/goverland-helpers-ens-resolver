# Goverland ENS Resolver Helper

<a href="https://github.com/goverland-labs/goverland-helpers-ens-resolver?tab=License-1-ov-file" rel="nofollow"><img src="https://img.shields.io/github/license/goverland-labs/goverland-helpers-ens-resolver" alt="GPL 3.0" style="max-width:100%;"></a>
![unit-tests](https://github.com/goverland-labs/goverland-helpers-ens-resolver/workflows/unit-tests/badge.svg)
![golangci-lint](https://github.com/goverland-labs/goverland-helpers-ens-resolver/workflows/golangci-lint/badge.svg)

## Overview

A gRPC microservice that resolves Ethereum Name Service (ENS) domains to addresses and vice versa. It acts as a centralized ENS resolution layer for the Goverland platform, abstracting away the complexity of interacting with multiple ENS data providers.

## What It Does

The service provides three core resolution capabilities:

1. **ENS Name → Address**: Resolve ENS domains (e.g., `vitalik.eth`) to Ethereum addresses
2. **Address → ENS Name**: Reverse resolve Ethereum addresses to their primary ENS names
3. **All ENS Names for Address**: Get all ENS domains owned by a specific address

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    goverland-helpers-ens-resolver               │
├─────────────────────────────────────────────────────────────────┤
│  gRPC Server (:20200)                                           │
│    └── EnsHandler                                               │
│          ├── ResolveAddresses (ENS name → address)              │
│          ├── ResolveDomains (address → ENS name)                │
│          └── ResolveAllDomains (address → all owned ENS names)  │
├─────────────────────────────────────────────────────────────────┤
│  Resolution Providers                                           │
│    ├── Stamp.fyi Client (primary resolver with caching)         │
│    │     └── In-memory cache for resolved addresses             │
│    └── Alchemy NFT API Client (for all domains lookup)          │
│          └── Queries ENS NFT contract ownership                 │
└─────────────────────────────────────────────────────────────────┘
```

## gRPC API

**Service Definition** (`protocol/enspb/ens.proto`):

```protobuf
service Ens {
  // Resolve ENS names to Ethereum addresses
  // Input: ["vitalik.eth", "ens.eth"] → Output: [{address, ens_name}, ...]
  rpc ResolveAddresses(ResolveAddressesRequest) returns (ResolveResponse);

  // Resolve Ethereum addresses to their primary ENS names
  // Input: ["0x123...", "0x456..."] → Output: [{address, ens_name}, ...]
  rpc ResolveDomains(ResolveDomainsRequest) returns (ResolveResponse);

  // Get ALL ENS names owned by an address (not just primary)
  // Input: "0x123..." → Output: ["name1.eth", "name2.eth", ...]
  rpc ResolveAllDomains(ResolveAllDomainsRequest) returns (ResolveAllDomainsResponse);
}
```

**Default Port**: `:20200`

## External Dependencies

| Provider | Purpose | Configuration |
|----------|---------|---------------|
| [Stamp.fyi](https://stamp.fyi) | Primary ENS resolution (address ↔ name) | `STAMP_ENDPOINT` |
| [Alchemy NFT API](https://docs.alchemy.com/reference/getnfts) | Query all ENS names owned by address | `ALCHEMY_API_KEY` |

### How Providers Are Used

**Stamp.fyi** (`pkg/sdk/stamp/`):
- Used for `ResolveAddresses` and `ResolveDomains`
- Calls `lookup_addresses` method to resolve addresses to ENS names
- Results are cached in-memory (no TTL, persistent for service lifetime)
- Batch processing: max 50 addresses per request

**Alchemy NFT API** (`pkg/sdk/alchemy/`):
- Used for `ResolveAllDomains`
- Queries the ENS NFT contract (`0x57f1887a8BF19b14fC0dF6Fd9B2acc9Af147eA85`) for all tokens owned by an address
- Extracts ENS names from NFT metadata

## How Stamp.fyi Resolves Names (Deep Dive)

Stamp.fyi is an open-source service by Snapshot Labs that resolves Web3 identities. The Goverland ENS resolver uses Stamp's JSON-RPC API.

### API Request Format

```bash
POST https://cdn.stamp.fyi/
Content-Type: application/json

{
  "method": "lookup_addresses",
  "params": ["0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"]
}
```

### Response Format

```json
{
  "result": {
    "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4": "aci.eth",
    "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045": "vitalik.eth"
  }
}
```

### Internal Resolution Flow

When Stamp receives a `lookup_addresses` request, it queries **multiple resolvers in parallel**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Stamp.fyi Server                                │
├─────────────────────────────────────────────────────────────────────────┤
│  POST / { method: "lookup_addresses", params: ["0x..."] }               │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Redis Cache Check                             │   │
│  │  Key: "address-resolvers:{address}"                              │   │
│  │  TTL: 43200 seconds (12 hours)                                   │   │
│  └──────────────────────────┬──────────────────────────────────────┘   │
│                              │ Cache miss                               │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │              Parallel Resolution (all run concurrently)          │   │
│  │                                                                   │   │
│  │  1. Snapshot Hub GraphQL                                         │   │
│  │     └── Query users by address for custom display names          │   │
│  │                                                                   │   │
│  │  2. ENS Reverse Records Contract ◄── PRIMARY FOR .eth NAMES      │   │
│  │     └── Contract: 0x3671aE578E63FdF66ad4F3E12CC0c0d71Ac7510C     │   │
│  │     └── Method: getNames(address[]) → string[]                   │   │
│  │     └── Returns primary ENS name set by user                     │   │
│  │                                                                   │   │
│  │  3. Unstoppable Domains API                                      │   │
│  │     └── Resolves .crypto, .nft, .x, .wallet, etc.                │   │
│  │                                                                   │   │
│  │  4. Lens Protocol GraphQL                                        │   │
│  │     └── Resolves .lens handles                                   │   │
│  │                                                                   │   │
│  │  5. Starknet ID                                                  │   │
│  │     └── Resolves .stark domains                                  │   │
│  │                                                                   │   │
│  │  6. Shibarium Name Service                                       │   │
│  │     └── Resolves Shibarium domains                               │   │
│  │                                                                   │   │
│  │  7. Space ID                                                     │   │
│  │     └── Resolves .bnb, .arb domains                              │   │
│  └──────────────────────────┬──────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Result Merging: First non-empty result wins (priority order)    │   │
│  │  → Cache result in Redis with 12-hour TTL                        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

### ENS Resolution Details

For `.eth` names specifically, Stamp calls the **ENS Reverse Records** contract on Ethereum mainnet:

```typescript
// Contract ABI
const abi = ['function getNames(address[] addresses) view returns (string[] r)'];

// Contract address (ENS Reverse Records)
const contract = '0x3671aE578E63FdF66ad4F3E12CC0c0d71Ac7510C';

// Call via Snapshot's brovider (RPC aggregator)
const names = await call(provider, abi, [contract, 'getNames', [addresses]]);
```

**Key points:**
- Returns the **primary ENS name** that the address owner has configured
- Uses the ENS "reverse resolution" mechanism (address → name)
- Names are validated with `ens_normalize()` to ensure they're valid ENS names
- Only returns names where the reverse record is properly set

### Name → Address Resolution (resolve_names)

For the reverse direction (ENS name to address), Stamp uses two strategies:

1. **ENS Subgraph** (primary): Query The Graph's ENS subgraph
   ```graphql
   query Domains($handles: [String!]!) {
     domains(where: {name_in: $handles}) {
       name
       resolvedAddress { id }
     }
   }
   ```

2. **Direct RPC** (fallback): If subgraph fails, call `provider.resolveName(handle)` for each unresolved name

### Caching Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                    Two-Level Caching                             │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Level 1: Stamp.fyi Redis                                        │
│  ├── Key: "address-resolvers:{address}"                          │
│  ├── TTL: 43200 seconds (12 hours)                               │
│  └── Shared across all Stamp instances                           │
│                                                                  │
│  Level 2: goverland-helpers-ens-resolver In-Memory               │
│  ├── Key: address (lowercase)                                    │
│  ├── TTL: None (service lifetime)                                │
│  └── Local to each resolver instance                             │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### Supported Domain Types

| Provider | Domains | Example |
|----------|---------|---------|
| ENS | `.eth` | `vitalik.eth` |
| Unstoppable Domains | `.crypto`, `.nft`, `.x`, `.wallet`, `.bitcoin`, `.dao`, `.888`, `.zil`, `.blockchain` | `brad.crypto` |
| Lens Protocol | `.lens` | `stani.lens` |
| Starknet ID | `.stark` | `checkpoint.stark` |
| Space ID | `.bnb`, `.arb` | `binance.bnb` |
| Shibarium | Shibarium domains | - |
| Snapshot | Custom display names | - |

## Configuration

Environment variables (see `.env.dist`):

```bash
# Logging
LOG_LEVEL=info              # zerolog level

# gRPC Server
GRPC_LISTEN=:20200          # gRPC server bind address

# Stamp.fyi (primary resolver)
STAMP_ENDPOINT=https://cdn.stamp.fyi

# Alchemy (for all-domains lookup)
ALCHEMY_API_KEY=            # Required for ResolveAllDomains

# Infrastructure (unused but configured)
INFURA_API_KEY=
INFURA_ENDPOINT=https://mainnet.infura.io/v3/

# Monitoring
PROMETHEUS_LISTEN=:2112
HEALTH_LISTEN=:3000
```

## Usage by Other Services

This service is consumed by:

- **goverland-core-storage**: Resolves ENS names for proposal authors and voters
- **goverland-inbox-storage**: Resolves ENS names for user profiles

Example gRPC client usage:

```go
import "github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"

conn, _ := grpc.NewClient("localhost:20200", grpc.WithInsecure())
client := enspb.NewEnsClient(conn)

// Resolve addresses to ENS names
resp, _ := client.ResolveDomains(ctx, &enspb.ResolveDomainsRequest{
    Addresses: []string{"0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4"},
})
// resp.Addresses[0].EnsName = "aci.eth"

// Resolve ENS names to addresses
resp, _ := client.ResolveAddresses(ctx, &enspb.ResolveAddressesRequest{
    Domains: []string{"vitalik.eth"},
})
// resp.Addresses[0].Address = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
```

## Build & Run

```bash
# Build
go build ./...

# Run tests
go test ./...

# Lint
golangci-lint run

# Run locally
cp .env.dist .env
# Edit .env with your API keys
go run main.go
```

## Directory Structure

```
internal/
├── app.go              # Application bootstrap
├── cache/              # In-memory cache for resolved names
├── config/             # Configuration structs
├── metrics/            # Prometheus metrics and request watcher
├── models/             # Domain models (ResolvedModel)
└── server/
    ├── handlers_ens.go # gRPC handler implementations
    └── forms/          # Request validation

pkg/
├── grpcsrv/            # gRPC server utilities
├── health/             # Health check server
├── prometheus/         # Metrics server
└── sdk/
    ├── alchemy/        # Alchemy NFT API client
    ├── infura/         # Infura client (currently unused)
    └── stamp/          # Stamp.fyi client with caching

protocol/
└── enspb/              # Generated protobuf code
    └── ens.proto       # Service definition
```

## Caching Behavior

The Stamp client maintains an in-memory cache:
- **Cache key**: Ethereum address (lowercase)
- **Cache value**: Primary ENS name
- **TTL**: None (cached for service lifetime)
- **Invalidation**: Service restart

This means ENS name changes won't be reflected until the service restarts or the address wasn't previously cached.

## Contribution Rules

[CONTRIBUTING.md](CONTRIBUTING.md)

## Changelog

[CHANGELOG.md](CHANGELOG.md)
