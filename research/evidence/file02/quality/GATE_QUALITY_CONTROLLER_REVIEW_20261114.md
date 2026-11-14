# Gate QUALITY — Controller technical review of Draft PR #28

Date: 2026-11-14. Review target: [PR #28](https://github.com/the-phong/Beam/pull/28),
Q evidence commit `f629b47c67dd75d929936a98aebf4617f513d108`, base
`test-mainet` merge commit `2b453b3c14c410545d5699e03cff90d0eb21d63e`.
**Technical recommendation: CONDITIONAL ACCEPT within the documented
single-host RPC scope, with five Q PASS and two Q PARTIAL.** This is a
Controller review recommendation; the Owner has not accepted Gate QUALITY
or authorized PR merge or X/D execution. This desk review was performed in
the same agent workspace as the engineering work, so it is not an
independent external reviewer sign-off for OUT-08.

## Review checks

| Check | Result |
|---|---|
| PR integration | Draft OPEN, mergeable to `test-mainet`; `scan`, `smoke`, two `verify` jobs all SUCCESS. |
| Source and plan lock | Q01–Q07 plans/audits and Q02 `schedule.v2` correction are in the source history. Corrected Q02 source commit `42d9ac9` precedes the new 96-network send. Historical Q02 v1 and abort are retained and excluded from current cross-profile inference. |
| Raw/plan integrity | The public [provenance index](DATA_PROVENANCE_INDEX.csv) checks local plan/manifest hashes against QC for 30 frozen sessions. Private raw remains ignored. The index marks two aborted sessions and the superseded 96-row Q02 v1 separately. |
| Independent Q verifiers | `python3 research/tools/file02_q01.py verify`; `file02_q02.py verify research/artifacts/file02/q02-20261114T044327Z-e7f236e2`; `file02_q03.py verify`; `file02_q04.py verify`; `file02_q06.py verify`; `file02_q07.py verify`: all exit 0. Q01: 76 pairs; Q02 v2: QC PASS; Q03: 22 groups all `INSUFFICIENT_SAMPLES`; Q04: 248 networks, zero exclusions; Q06: `PARTIAL_HOST_REPLICATION`; Q07: 248 valid, zero exclusions. |
| Correctness and statistics | Q07's final full-genesis replay addresses the earlier unused-sender audit assumption. Run/network or fixture block is the CI unit; a transaction within a run is not a replicate. Null inferential p99 stays null. |

## Disposition and conditions

Q01–Q04 and Q07 have sufficient bounded evidence to retain PASS and checked
boxes. Q05 cannot identify a saturated bottleneck within 25 tx/s; stage
crypto CPU and full queue observability are absent. Q06 has seed/order
replication but no independent physical-host replication. Their boxes must
remain open as PARTIAL. Neither a `PASS` technical QC session nor a manifest
hash alone establishes inferential validity: Q02 v1 is the concrete example.

The RQ1b finding remains **no observed loss at the preregistered 3-second
fixed-legitimate first-all-four observation endpoint in tested 20-tx/s
matched single-host RPC strata**. Its [0,0] paired intervals are degenerate
because all runs reached the endpoint ceiling; they do not establish
equivalence, exact commit latency or a universal no-effect conclusion. Q02
compares structural-invalid traffic only; C/P wire sizes are not matched to
H, while CP/H are matched. No natural knee, multi-host effect, P2P-origin
effect or DoS mechanism is inferred.

The review found no reason to convert Q05/Q06 to PASS or to rerun workloads
for a positive result. A prospective mechanism/service extension may be
proposed separately in the [extension proposal](../../../protocol/FILE02_MECHANISM_SERVICE_EXTENSION_PROPOSAL.md), retaining the accepted RQ1 and RQ1b definitions and
the historical 3-second result. It must lock new estimands, precision,
instrumentation, safety and QC before any new run. D01–D13 remain NOT_RUN
until an explicitly accepted target and appropriate event/label exist.

**Owner action needed:** decide whether to accept Gate QUALITY conditionally
on the above boundaries and, only after acceptance, whether to merge PR #28.
The Owner's decision is deliberately left `PENDING` in
[GATE_QUALITY_OWNER_DECISION.md](GATE_QUALITY_OWNER_DECISION.md). This review
does not merge the PR, modify raw data or open X/D execution.
