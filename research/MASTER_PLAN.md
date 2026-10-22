# Pluto research master plan

Project question: What are the resource costs and correctness properties of
hybrid secp256k1–ML-DSA-44 transactions in a custom CometBFT-based EVM
application under deliberately invalid workloads? No finding or publication
outcome is presumed. Every phase requires a separate human authorization;
finishing one phase does not authorize the next.

The method is: protocol specification → correctness oracle → functional
implementation → adversarial testing → instrumentation → baseline measurement
→ root-cause analysis → conditional mechanism design → comparative evaluation
→ reproducible artifact. Each phase records a report, exact revision, executed
checks, limitations and a review decision.

For locally developed Phase 02, the user explicitly replaced the Git commit
checkpoint with a SHA-256 Phase 01 baseline manifest. `BASE_COMMIT` and
`RESULT_COMMIT` are therefore not applicable; this does not authorize Phase 03
or relax any blockchain/correctness gate. Future phases should record a
reproducible source snapshot or checksum manifest when no commit is used.

| Phase | Objective | Main deliverables | Dependencies | Acceptance criteria | Scope restrictions |
|---|---|---|---|---|---|
| **01 — Project Foundation & Architecture Bootstrap** | Establish a reproducible, honest foundation. | Go module, metadata CLI, workflow, architecture, decisions, manifest, report. | None. | Build, tests, vet, formatting, CLI, documents and G1–G12 pass with evidence. | No blockchain, ABCI, EVM, hybrid crypto, admission control or benchmarks. |
| **02 — CometBFT Foundation** | Run a minimal real ABCI++ application against a pinned CometBFT version. | Node startup/configuration, lifecycle integration, durable minimal state and deterministic AppHash contract, replay tests. | Approved 01; final CometBFT version and API review. | Real single-node smoke test, restart/replay and state hash agreement at the tested scope. | No EVM or hybrid transaction claims. Do not alter consensus algorithms. |
| **03 — Classical Transaction Core** | Specify and implement canonical classical transaction handling. | Versioned transaction bytes, secp256k1 verification, account/nonce/balance rules, positive and negative vectors. | Approved 02; transaction protocol decision. | Exact signed bytes, binding, replay rejection and deterministic validity verified. | No ML-DSA or hybrid fallback policy yet. |
| **04 — EVM Integration** | Execute classical transactions sequentially with geth components. | EVM adapter, gas/receipt/state mapping, persistence and differential correctness tests. | Approved 03; geth and storage version decisions. | Deterministic execution, crash/replay safety and selected geth semantic comparisons pass. | No parallel executor or performance optimization. |
| **05 — Hybrid Protocol Design** | Freeze the hybrid security profile before code. | Canonical envelope spec, key policy, signed-intent domain, state rules, rejection oracle and test-vector design. | Approved 04; FIPS 204 and library review. | Independent review finds unambiguous encoding, same-intent authentication and fail-closed policy. | Specification only; no implementation or admission mechanism. |
| **06 — Hybrid Cryptographic Module** | Implement and validate both proof checks and binding primitives. | Library adapters, constrained parsing, vectors, negative tests and crypto review record. | Approved 05; pinned crypto libraries. | Known-answer/interoperability vectors and malformed-input tests pass. | No custom cryptographic primitive and no transaction integration claim. |
| **07 — Hybrid Transaction Integration** | Apply hybrid policy to application transactions. | End-to-end account binding, canonical decode, state validation, EVM mapping and fail-closed tests. | Approved 06. | Hybrid-valid tx executes; either missing/invalid proof rejects; state transitions replay identically. | No optimization or benchmark conclusion. |
| **08 — ABCI Lifecycle & Safety** | Validate the hybrid policy across real ABCI++ paths. | CheckTx/recheck, proposal and finalize invariants; adversarial proposal tests; AppHash consistency. | Approved 07. | Valid proposals remain acceptable; malicious input cannot bypass checks or split deterministic state. | Do not treat local admission as consensus proof. |
| **09 — Adversarial Testing** | Build an oracle-backed invalid-input corpus. | Malformed, one-proof-invalid, binding, state, replay and proposer-injection cases; fuzz/property tests. | Approved 08. | Every case has input bytes, expected path/outcome and reproducible seed; no safety regression. | No performance or novelty claim from functional tests alone. |
| **10 — Benchmark Infrastructure** | Measure costs by phase without changing validity. | Workload generator, instrumentation, profiles, environment/seed manifests and clean controls. | Approved 09. | Repeatable runs capture CPU, allocations, latency, throughput and resource context with validation. | No admission-control optimization. |
| **11 — Baseline Evaluation** | Quantify and diagnose invalid-input effects. | Single-node and appropriate baseline measurements, statistical analysis, profiles and root-cause report. | Approved 10. | Raw data, scripts, controls, uncertainty and interpretation are independently reviewable. | No invented improvement or assumption that a mechanism is needed. |
| **12 — Research Mechanism P (conditional)** | If 11 identifies an addressable bottleneck, design and test a targeted mechanism. | Explicit go/no-go decision, safety proof obligations, implementation, ablation and comparison. | Approved 11 **plus** documented research decision that a mechanism is justified. | If authorized: correctness preserved and measured benefit/trade-offs shown against strong baselines; if not: record no-go and skip. | Never implement P merely to follow the roadmap; no consensus-critical local nondeterminism. |
| **13 — Multi-Node Experiments** | Test correctness and costs across validators and hosts. | Repeatable local and multi-host deployments, workload runs, AppHash agreement and network-resource observations. | Approved 11, and 12 only if pursued. | Independent validators agree on state; environment, topology and raw results are recorded. | Local multi-process tests must not be misrepresented as multi-host results. |
| **14 — Research Artifact & Manuscript** | Package a reproducible, evidence-limited research contribution. | Pinned source/data, scripts, artifact instructions, threat model, limitations and manuscript. | Approved 13 and complete review of prior phases. | Independent reproduction of supported claims and clear distinction between observation and hypothesis. | No unsupported novelty, vulnerability or venue acceptance claim. |

If Phase 11 finds no meaningful addressable cost, Phase 12 is skipped or
reframed only after a documented decision. The research can still contribute
a rigorous correctness and cost characterization. Gate decisions and accepted
changes are recorded in `DECISIONS.md` and `PHASE_STATUS.md`.
