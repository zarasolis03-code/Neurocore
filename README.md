
# Neurocore (prototype)

Neurocore is a minimal Go-based blockchain prototype featuring an AI-Workload PoW simulation and an adaptive difficulty algorithm. This repository is intended as a starting point.

Quick run:

```bash
go run .
```

The program will create a `neurocore_history.json` file with the genesis block, start a miner goroutine and a TCP P2P server on port 40404, and print a simple dashboard.

Quick Start (Raspberry Pi / Linux)

1. Build and run:

```bash
cd /workspaces/Neurocore
go build -o neurocore
./neurocore
```

2. Run using `go run` (dev):

```bash
go run .
```

3. Inspect chain or submit a sample transaction:

```bash
# print chain
go run . -inspect

# produce a sample tx JSON
go run . -submit
```

Repository structure

- `blockchain.go` — core chain, adaptive difficulty, persistence
- `wallet.go` — ECDSA wallet and helpers
- `miner.go` — AI-workload PoW miner simulation
- `network.go` — minimal TCP P2P server (port 40404)
- `main.go` — startup, wallet creation, miner and P2P runner, dashboard
- `cli.go` — small CLI helpers (`-inspect`, `-submit`)
- `neurocore_history.json` — chain history (created at runtime)

License

This project is released under the MIT License. See `LICENSE` for details.

# Neurocore