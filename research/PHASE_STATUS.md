# Phase status — 2026-10-23

## Current Phase 03C implementation candidate

The Owner authorized Phase 03C after the Phase 03B governance closeout merged
into `test-mainet` at `61bc65fdebf358fb0330cb9ebd4a332ca4a06aca`.
Work on `phase/03c-b0-state-persistence` adds standalone V2 account state,
genesis parsing, sequential simulation, AppHash and a separate Pebble store.
The Phase 02 ABCI runtime still rejects every transaction and uses its V1
record and AppHash. D-020–D-023 remain accepted and unchanged. D-024 proposes
the physical V2 storage layout for Controller and Owner review; it is not yet
accepted. Phase 03C completion is subject to evidence, CI and independent
review. Phase 03D and all later phases are **NOT AUTHORIZED**.

## Historical governance closeout after Phase 03B merge

Phase 03B is **completed, Controller-accepted, and merged** via PR #4 into
`test-mainet` at `8be07c87abed73763d40252e99bbf0573da487b4`. The Owner
accepted D-023 on 2026-10-23; D-020–D-022 remain accepted and unchanged.
Both post-merge GitHub verification runs passed. The only active authorization
is documentation-only acceptance closeout on
`governance/03b-acceptance-closeout`, with a Draft PR and Owner-controlled
merge. Phase 03C is **NOT AUTHORIZED**. The running node remains the Phase 02
transaction-free application; `internal/tx` is a standalone stateless library.
Account state, AppHash V2 and ABCI transaction handling remain unimplemented.

## Historical Phase 03B execution status (superseded by Controller acceptance and merge)

The Owner accepted D-020–D-022 and merged Phase 03A PR #3 into `test-mainet`
at `54bf720ac4625c17eaa6101a8f398ef7c0663432`. The Owner separately
authorized Phase 03B on `phase/03b-classical-tx-crypto`. This phase implements
and verifies a standalone stateless transaction library and public test
vectors. The Phase 02 node still rejects all transactions; account state,
AppHash V2 and ABCI integration are reserved for Phases 03C–03E. Phase 03C is
not authorized. Completion and acceptance remain subject to Controller/Owner
review of the Phase 03B report and Draft PR.

## Historical Phase 03A closeout status (superseded by merger and 03B authorization)

- **Active authorization:** Phase 03A documentation closeout only. The Owner
  authorized Phase 03A after Phase 02R was merged into `test-mainet` at
  `ab1d1746fdd70bc0bbf00f20a186663f3cacf2ab`, then accepted the B0 V1
  protocol and D-020–D-022 after Controller R1–R5 rework on 2026-10-23.
- **Engineering status:** Four B0 protocol documents and the Phase 03A report
  specify canonical transaction bytes, secp256k1 signed intent, account/state
  transitions, AppHash V2, ABCI validity and negative-vector plans. The Owner
  accepted the normative rules after the Controller's R1–R5 rework; they are
  **accepted as protocol, not implemented**. PR #3 remains a Draft awaiting
  owner-controlled merge. The runtime remains the Phase 02 empty-block
  application.
- **Verification:** `git diff --check`, `make verify` and
  `make test-integration` exited 0 on the Phase 03A branch; exact commands,
  output, toolchain and limits are in `research/evidence/phase_03a/` and
  `research/PHASE_REPORTS/PHASE_03A_REPORT.md`. The same required checks and
  a runtime/dependency/test source-diff check exited 0 during controller
  rework; see `controller_rework_verification.txt`. Documentation closeout
  checks also exited 0; see `closeout_verification.txt`. No cryptographic
  vector was generated or executed in this specification phase.
- **Historical baseline:** Phase 02/02R architecture, reports and evidence
  remain preserved. Phase 02 V1 AppHash and `PLUTO02!` storage are still the
  only implemented state format. Accepted V2 requires a fresh B0 chain; no
  automatic migration is specified or implemented.
- **Governance:** D-018 branch/PR/owner-merge control applies. D-019 requires
  B0-Core before the conditional Phase 04 EVM decision. D-020–D-022 are
  **Accepted by Owner** as protocol requirements for later authorized work.
- **Known limits:** Current tests exercise only the Phase 02 runtime. No B0
  transaction execution, account persistence, secp256k1 verification, V2
  commitment, four-validator agreement or security/performance result exists.
- **Next action:** Await the Owner-controlled merge of Draft PR #3 and a
  separate Phase 03B decision after documentation closeout.
  **Phase 03B is not authorized.** No merge by the agent.
