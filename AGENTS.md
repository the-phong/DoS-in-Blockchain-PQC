# Pluto agent execution contract

## Current Phase 02 review handoff (supersedes historical phase status below)

Phase 02 — CometBFT Foundation is authorized and implemented for independent
review. The repository now contains a real single-validator, empty-block
CometBFT v0.38.26 chain connected to Pluto ABCI++ and Pebble v1.1.5. Phase 03
is **not authorized**. An agent may work only within an explicitly authorized
phase; the Phase 02 report is not approval for Phase 03. No transaction
execution, EVM, hybrid authentication or admission optimization exists.

The user's local-only decision supersedes the original Phase 02 Git checkpoint:
do not require or create commits, branches, remotes, Git identity or pushes.
Phase 01 baseline checksums are in `research/evidence/phase_02/` and the
Phase 01 report/evidence must remain preserved. Record further changes in
research governance documents and evidence, not by claiming a Git revision.

For Phase 02 verification run `make verify`, `make test-integration`, and the
relevant race test. The application is transaction-free and all unsupported
transactions must fail closed. `FinalizeBlock` is pending-only; `Commit`
performs the synced persistent write. The Phase 02 AppHash is a fixed digest
of version, chain ID and height, not an Ethereum state root. Generated
validator/node keys must never be included in reports or commits. Do not
broaden node exposure beyond localhost without a new approved decision.

Pluto is a greenfield research project studying resource cost and processing
correctness for hybrid secp256k1–ML-DSA-44 transactions under deliberately
invalid workloads in a custom CometBFT-based EVM application. The repository
currently contains a Phase 02 transaction-free single-validator chain, not an
EVM or hybrid blockchain.

## Phase authority

Active authorized phase: **02 — CometBFT Foundation**, for review and any
authorized Phase 02 rework only. An agent may work only within the explicitly
authorized phase. Phase 03 remains unapproved. Necessary work outside
the active phase must be documented and reported, never silently implemented.
Do not infer approval from a PASS verdict or from roadmap text.

## Architecture and security rules

- CometBFT owns consensus, P2P, block coordination and its mempool. Pluto owns
  application transactions, authentication, account policy, state, EVM
  execution and AppHash. Do not modify CometBFT consensus algorithms.
- The target is a custom Go EVM application on CometBFT. Do not add Cosmos SDK
  without an explicit recorded decision.
- Preserve sequential EVM execution as the initial baseline. Do not design an
  admission optimization before Phase 11 profiling justifies it.
- Consensus-critical results must be deterministic for the same relevant
  committed state and proposal. Never use local clocks, queue size, CPU
  budgets, nonreplicated caches, scheduling or mempool arrival order as
  consensus validity inputs.
- Under future `HYBRID_REQUIRED`, both secp256k1 and ML-DSA-44 must verify the
  same canonical intent and bind to registered keys. Fail closed; no fallback.
  Do not invent cryptographic primitives.
- `CheckTx` is an admission filter; proposal validation must independently
  handle transactions that did not pass a local `CheckTx`.
- Do not equate AppHash with an Ethereum state root without a documented
  commitment design and tests.

## Implementation conventions

- Use idiomatic Go, standard library where practical, explicit errors and
  small packages with a current purpose. Avoid speculative interfaces and
  empty future modules.
- Pin exact third-party module versions when first imported. Record source,
  minimum Go version, reason, risks and actual build/test evidence in
  `research/DECISIONS.md` and the manifest. Never silently upgrade an
  accepted version. `go.mod` is the executable dependency source of truth.
- Add meaningful tests for behavior implemented now, including invalid
  inputs and deterministic replay when those become relevant. Run `make verify`
  before declaring a phase done; report every failure and skip.
- Do not claim a blockchain, EVM, hybrid verifier, benchmark result or
  multi-node experiment exists until its implementation and tests are present.

## Documentation and evidence

- Keep `research/ARCHITECTURE.md`, `research/MASTER_PLAN.md`,
  `research/DECISIONS.md`, `research/PHASE_STATUS.md`, and
  `research/config/manifest.yaml` aligned with actual code and approvals.
- Put exact commands, exit codes, relevant output and limitations in a phase
  report and `research/evidence/` as appropriate. Record revision, toolchain,
  workload seeds, hardware and configuration for later experiments.
- A phase report is an engineering claim for independent review, not approval.
  Stop at the end of Phase 02 after its report. Record the reviewer's decision
  before starting Phase 03.
