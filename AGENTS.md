# Pluto agent execution contract

## Current authority — Phase 03A (2026-10-23)

The project owner explicitly authorized **Phase 03A — Classical Transaction
Protocol Specification** after the Phase 02R merge and has now accepted its
B0 V1 protocol and D-020–D-022 following Controller rework. The coding agent
is authorized for Phase 03A documentation closeout only on the existing
`phase/03a-tx-protocol` branch from verified `test-mainet`. It may commit and
push closeout documentation to the existing Draft PR targeting `test-mainet`.
Owner acceptance freezes the specified rules for later separately authorized
work; it does not mean those rules are implemented or tested in the node.
Runtime changes, transaction/crypto implementation, EVM, hybrid mechanisms,
consensus/mempool changes and Phase 03B or later work are **not authorized**.
The implemented Phase 02 architecture remains the transaction-free baseline.

ChatGPT acts as technical controller; the Coding Agent performs engineering
work. The project owner approves phase transitions and performs or explicitly
approves merges. Agent-controller technical communication, source documents
and reports are in English. The controller provides explanations and decision
requests to the owner in Vietnamese.

Every authorized phase uses its own branch from the verified `test-mainet`
baseline. Normal commits and pushes are permitted within that phase's scope;
open a Draft PR against `test-mainet`. Never push directly to `test-mainet` or
`master`, merge automatically, force-push, or start another phase without
separate approval. Preserve unrelated local work. End every execution with a
full report including source/result commit SHAs, branch, PR URL or blocker,
verification and limitations, then STOP. A passing report never grants the
next phase.

The owner approved Master Roadmap v2: finish minimal B0-Core through Phase
03E before the Phase 04 decision on whether sequential EVM integration is
necessary. EVM is optional pending that decision; if approved, execution must
begin sequentially. Phase 03A documentation closeout alone is now authorized;
03B–14 remain future work until separately authorized. Do not design admission
optimization before Phase 11 profiling.

## Historical Phase 02R authority (superseded by owner Phase 03A decision)

Phase 02R was authorized for independent audit, verification, governance
reconciliation and a Draft PR on `phase/02r-audit-governance`. Its report did
not itself authorize Phase 03A. At that time, Phase 03A implementation was
unapproved. The subsequent explicit owner decision above supersedes that
status for specification work only; the Phase 02R report and evidence remain
preserved.

## Historical Phase 02 review handoff (superseded by Phase 02R authority)

Phase 02 — CometBFT Foundation is authorized and implemented for independent
review. The repository now contains a real single-validator, empty-block
CometBFT v0.38.26 chain connected to Pluto ABCI++ and Pebble v1.1.5. Phase 03
is **not authorized**. An agent may work only within an explicitly authorized
phase; the Phase 02 report is not approval for Phase 03. No transaction
execution, EVM, hybrid authentication or admission optimization exists.

For Phase 02, the user's local-only decision superseded its original Git
checkpoint: no commits, branches, remotes, Git identity or pushes were
required or created during that phase. D-018 supersedes this policy for newly
authorized phases while preserving that historical fact.
Phase 01 baseline checksums are in `research/evidence/phase_02/` and the
Phase 01 report/evidence must remain preserved. The then-current instruction
recorded further changes in research governance documents and evidence without
claiming a Git revision; D-018 supersedes that instruction for new phases.

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

## Historical Phase 02 authority (superseded)

The then-active authorized phase was **02 — CometBFT Foundation**, for review
and authorized Phase 02 rework only. An agent may work only within the
explicitly authorized phase. Phase 03 remains unapproved. Necessary work outside
the active phase must be documented and reported, never silently implemented.
Do not infer approval from a PASS verdict or from roadmap text.

## Architecture and security rules

- CometBFT owns consensus, P2P, block coordination and its mempool. Pluto owns
  application transactions, authentication, account policy, state, EVM
  execution and AppHash. Do not modify CometBFT consensus algorithms.
- The historical target was a custom Go EVM application on CometBFT. D-019
  makes EVM conditional at Phase 04. Do not add Cosmos SDK without an
  explicit recorded decision.
- If EVM is approved, preserve sequential execution as its initial baseline.
  Do not design an admission optimization before Phase 11 profiling justifies it.
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
  Stop at the end of Phase 03A after its report. Record the controller's and
  owner's decision before starting Phase 03B.
