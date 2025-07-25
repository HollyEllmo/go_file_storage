# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

- `make build` - Builds the application to `bin/fs`
- `make run` - Builds and runs the application
- `make test` - Runs all tests with verbose output (`go test ./... -v`)

## Code Architecture

This is a distributed file storage system implemented in Go with P2P networking capabilities. The system allows nodes to store, retrieve, and replicate files across a network.

### Core Components

**FileServer** (`server.go`): The main orchestrator that manages file operations and network communication. It handles:
- File storage and retrieval operations
- Broadcasting messages to peers for file operations
- Managing peer connections and network topology
- Coordinating encrypted file transfers

**Store** (`store.go`): File system abstraction that handles local disk operations using content-addressed storage (CAS). Key features:
- Path transformation using SHA1 hashing for content addressing
- Support for encrypted file storage/retrieval
- Organized storage by node ID to enable multi-tenant file systems

**P2P Transport Layer** (`p2p/`): Network communication infrastructure including:
- `Transport` interface for pluggable network protocols
- `TCPTransport` implementation for TCP-based communication
- `Peer` abstraction for remote node connections
- Message encoding/decoding with support for both regular messages and file streams

**Crypto** (`crypto.go`): Encryption utilities providing AES-CTR encryption for file content with random IV generation.

### Message Types and Network Protocol

The system uses two main message types:
- `MessageStoreFile`: Broadcasts file availability to network peers
- `MessageGetFile`: Requests file retrieval from network peers

Network communication supports two modes:
- Regular messages (using `IncomingMessage` byte marker)
- File streams (using `IncomingStream` byte marker)

### Key Architecture Patterns

- **Content-Addressed Storage**: Files are stored using SHA1 hash-based paths for deduplication
- **Peer-to-Peer Replication**: Files are automatically replicated across connected nodes
- **Transport Abstraction**: Network layer is abstracted to support different protocols
- **Encrypted Storage**: All files are encrypted before network transmission using node-specific keys

### Entry Point

`main.go` demonstrates the system by:
1. Creating 3 server nodes on ports 3000, 4000, and 6000
2. Establishing a bootstrap network topology
3. Storing test files and demonstrating retrieval after local deletion