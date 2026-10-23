# Phase status — 2026-10-23

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
