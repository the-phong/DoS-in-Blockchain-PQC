# Proposed File 02 mechanism and service measurement addendum

**Status: PROPOSED FOR OWNER REVIEW — NOT LOCKED, NOT RUN-AUTHORIZED.**
Date: 2026-11-14. This document records the proposed X01–X18 extension
without amending accepted RQ1 (H−CP application CheckTx) or the separate
RQ1b 3-second legitimate-service endpoint. It does not open D01–D13,
authorize ML, alter S06 caps, or imply that Gate QUALITY is accepted.

## Scientific questions and preserved results

- RQ1 remains the accepted controlled resource-cost comparison. The
  historical H−CP result is +81.47 µs/tx, paired 95% CI
  [+75.16, +87.26] µs/tx at 20 valid tx/s in the tested single-host RPC
  stratum. It is application CheckTx wall/service time, not isolated
  ML-DSA CPU.
- RQ1b remains the accepted fixed-legitimate H stream compared under
  matched benign and invalid background. Its 3-second first-all-four
  observation result is negative and ceiling-limited in the tested strata.
  A new continuous-latency estimand must be explicitly prospective for
  confirmatory use; historical reanalysis is exploratory only.
- Proposed RQ-Mechanism asks where invalid-workload verification work is
  spent, why the registered service endpoint did not degrade, and whether
  measured queue/resource accumulation occurs within the authorized safe
  range. It is a new question requiring an Owner decision and a locked
  design before data collection.

## Proposed ordered work, without implicit authority

| Step | Proposed outputs | Required lock before any new workload |
|---|---|---|
| X01–X04 | Causal map; historical QC on original tool commits; exploratory existing-latency note; separate estimands and practical margins. | Declare continuous-latency estimand, precision and decision rule without changing historical 3-second result. |
| X05–X09 | Dedicated app-process placement if feasible; validated CPU attribution; queue semantics; precise observation contract; instrumentation ON/OFF and correctness QA. | Show measurement overhead, process identity, clock uncertainty, unchanged validation/state and safety abort behavior. Keep stage CPU `null` until independently validated. |
| X10–X14 | Prospective three-arm L/B/I controls, early/late-invalid and replay controls, safe-range sensitivity, unavailable-source matrix. | Fresh matched networks, fixed legitimate identities/offsets/genesis, balanced arm order, bounded run-level n and rates. Maintain ≤25 tx/s, ≤16 outstanding, ≤60 s/run and all other S06 caps. No P2P-origin/multi-host claim without separate capability and authority. |
| X15–X18 | Locked inference/precision plan; mechanism analysis; independently verified immutable package; service/mechanism gate decision. | Run/network is the independent unit, per-contrast CI and missingness rules; equivalence claims require a prospective margin and new precision-adequate data. Retain failed/excluded sessions and decide negative/positive/unsupported outcomes without post-hoc thresholding. |

The distinct L→I (adding invalid background) and B→I (replacing benign
background at matched offered rate/bytes) contrasts must not be pooled or
treated as the same causal effect. Stage wall spans, process CPU and
per-stage CPU remain different quantities. Socket CheckTx wait is not a
full transaction queue; Comet relay envelope bytes are not per-peer
latency; a 100-ms block observation is an inclusion upper bound, not exact
commit time. The absence of a knee ≤25 tx/s is `NOT_OBSERVED_WITHIN_SAFE_RANGE`,
never evidence that 25 tx/s is capacity.

## Detection gate stays separate

Original early service-degradation RQ2 requires a prospectively defined,
measured degradation event and features available before that event. With
no event, its status is `NOT_ESTABLISHED` and lead time is `null`, not zero.
An alternative adversarial-workload detector needs its own Owner-approved
question, labels and natural-rejection controls. RQ3 requires independent
holdout comparison with threshold, EWMA/CUSUM and cost-aware rules at a
common FPR, lead-time and overhead budget. This proposal neither trains ML
nor claims a detector is needed.

## Decision required to lock the extension

After bounded Gate QUALITY review, the Owner must explicitly accept or
revise the proposed mechanism question, permitted measurement architecture,
case set, topology, sample-size/precision plan, safety and release branch.
Only then may a separate locked extension protocol authorize new X runs.
This proposal is not that authorization. The final integrated report and
OUT-01–OUT-08 remain future deliverables with unsupported/NOT_RUN sections
reported honestly rather than fabricated.
