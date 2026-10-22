# Phase status — 2026-10-23

- **Active authorization:** Phase 03A — classical transaction protocol
  specification only. The project owner explicitly authorized this phase
  after Phase 02R was merged into `test-mainet` at
  `ab1d1746fdd70bc0bbf00f20a186663f3cacf2ab`.
- **Engineering status:** Four B0 protocol documents and a Phase 03A report
  propose canonical transaction bytes, secp256k1 signed intent, account/state
  transitions, AppHash V2, ABCI validity and negative-vector plans. These
  rules are **proposed pending controller and owner review**, not accepted or
  implemented. The runtime remains the Phase 02 empty-block application.
- **Verification:** `git diff --check`, `make verify` and
  `make test-integration` exited 0 on the Phase 03A branch; exact commands,
  output, toolchain and limits are in `research/evidence/phase_03a/` and
  `research/PHASE_REPORTS/PHASE_03A_REPORT.md`. No cryptographic vector was
  generated or executed in this specification phase.
- **Historical baseline:** Phase 02/02R architecture, reports and evidence
  remain preserved. Phase 02 V1 AppHash and `PLUTO02!` storage are still the
  only implemented state format. V2 is a fresh-chain proposal with no
  automatic migration.
- **Governance:** D-018 branch/PR/owner-merge control applies. D-019 requires
  B0-Core before the conditional Phase 04 EVM decision. D-020–D-022 are
  proposals awaiting review; the owner has not accepted their normative rules.
- **Known limits:** Current tests exercise only the Phase 02 runtime. No B0
  transaction execution, account persistence, secp256k1 verification, V2
  commitment, four-validator agreement or security/performance result exists.
- **Next action:** Controller and owner review the Phase 03A protocol and
  Draft PR, resolve recorded questions and decide whether to accept the
  specification. **Phase 03B is not authorized.** No merge by the agent.
