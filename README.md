# Pluto

Pluto Phase 02 is a real, local, **single-validator CometBFT v0.38.26** chain
with a custom Go ABCI++ application, deterministic height-based AppHash and
durable Pebble v1.1.5 committed state. It produces empty blocks only; all
application transactions are rejected. This is a research foundation, not a
production chain, EVM, hybrid verifier or multi-validator experiment.

## Prerequisites and build

On Linux/WSL2, install Go 1.24.x or newer and GNU Make. Network access is
needed on first build to download exact pinned Go modules. Docker is not
required. From this directory:

```bash
make build
./build/plutod version
./build/cometbft version
```

`make build` creates **both** executables together under ignored `build/`.
`plutod` requires its adjacent pinned `cometbft` executable for `init` and
`start`; do not replace it with an arbitrary system-installed binary.

## Initialize, run, inspect and restart

Choose a private node home; the default `.pluto` is ignored by Git. Run from
one terminal:

```bash
./build/plutod init --home .pluto --chain-id pluto-dev-1
./build/plutod start --home .pluto
```

`init` creates CometBFT genesis, validator/node keys, local ports and Pluto
metadata. It refuses an existing home and never resets one. The default RPC,
P2P and ABCI ports are 26657, 26656 and 26658, bound to `127.0.0.1`; override
all three at init with `--rpc-port`, `--p2p-port`, `--abci-port` for another
local test. Empty blocks are configured every second. Protect files under
`.pluto/config/`, especially private validator and node keys; never publish
them or expose the development RPC/validator publicly.
Protect `config.toml` as well: startup does not recheck its listener and proxy
addresses against the localhost metadata written by `init`.

From another terminal, inspect real CometBFT RPC responses:

```bash
curl -fsS http://127.0.0.1:26657/status
curl -fsS 'http://127.0.0.1:26657/block?height=2'
curl -fsS http://127.0.0.1:26657/abci_info
```

Wait for at least height 2. `/status` gives the consensus height; `/abci_info`
gives the latest committed Pluto state height and base64 AppHash. A block
header's `app_hash` corresponds to the **previous** committed state height,
not the same-numbered `/abci_info` height. Stop with Ctrl-C; then inspect the
durable Pluto record and restart from the **same home**:

```bash
./build/plutod state --home .pluto
./build/plutod start --home .pluto
```

Do not run `state` while the node is active: Pebble enforces exclusive writable
access. On restart, CometBFT logs its ABCI handshake/replay and continues block
production. If `init` fails because a home exists, use the existing home with
`start` or select a *new* directory; never delete a node directory just to
make the command pass. If startup fails, check local port collisions, the
adjacent CometBFT binary, genesis/Pluto metadata mismatch, database lock and
key-file permissions. Partial initialization is not automatically cleaned up.

## Verification and research records

```bash
go build ./...
go test ./...
go vet ./...
make verify
make test-integration
go test -race ./internal/app ./internal/store
```

`make verify` builds both executables, runs unit/ABCI tests, vet, a read-only
format check, module checksum verification and CLI version check.
`make test-integration` runs a bounded real-node test in a dedicated temporary
home: init, RPC-observed blocks, graceful stop, independent state reload,
restart and further blocks. It only cleans up its own processes/directories.
Actual Phase 02 results and limitations are in
[the Phase 02 report](research/PHASE_REPORTS/PHASE_02_REPORT.md), with local
Phase 01 baseline checksums in
[BASELINE_SHA256.txt](research/evidence/phase_02/BASELINE_SHA256.txt).
[Architecture](research/ARCHITECTURE.md), [decisions](research/DECISIONS.md),
[phase status](research/PHASE_STATUS.md) and the
[manifest](research/config/manifest.yaml) are the current governance records.
The [Phase 02R independent review](research/PHASE_REPORTS/PHASE_02R_REPORT.md)
records fresh checks and limitations. The 14-phase [master plan](research/MASTER_PLAN.md)
now places minimal B0-Core before a conditional Phase 04 EVM decision.
Transaction implementation in Phase 03B and later is **not authorized**. D-018 governs new work:
one approved phase per branch, Draft PRs to `test-mainet`, and owner-controlled
merges. The Phase 02 local-only Git exception remains in historical records.

## Phase 03A protocol proposal

The owner has authorized a specification-only Phase 03A. The proposed
[B0 transaction format](research/specs/B0_CLASSICAL_TRANSACTION_PROTOCOL_V1.md),
[account transitions and AppHash V2](research/specs/B0_STATE_TRANSITION_SPEC_V1.md),
[ABCI validity policy](research/specs/B0_ABCI_VALIDITY_SPEC_V1.md) and
[test-vector plan](research/specs/B0_TEST_VECTOR_PLAN_V1.md) are awaiting
controller and owner review. They are not implemented by the current binary;
it still rejects all transactions. Phase 03B is not authorized.

