# B0 V1 cryptographic and negative vector plan — Phase 03A proposal

**Status: PLAN ONLY.** No transaction bytes, signatures, addresses, hashes or test outcomes in this document have been generated or executed. The table gives expected results **if** the proposed protocol is approved and correctly implemented. Phase 03B must independently generate byte vectors with at least two implementations or a reviewed reference implementation, save exact hex input/output and tool/library versions, and run them. Phase 03C–03E must add state/ABCI/runtime evidence. Never use the keys below outside nonproduction tests.

## 1. Deterministic fixture recipe

- Chain ID: exact ASCII `pluto-b0-v1`. A separate chain fixture uses ASCII `pluto-other-1`.
- Labeled test-only secret scalars: `TEST_ONLY_A = 1` and `TEST_ONLY_B = 2`, each encoded as a 32-byte big-endian scalar for secp256k1 operations. Derive compressed public keys `P_A = 1*G`, `P_B = 2*G`, then addresses `A` and `B` using the proposed address domain. **These are public, insecure test scalars, not node or validator keys.** The spec intentionally does not assert uncomputed address hex.
- Fresh genesis state `S0`: one allocation `(A, 100 B0 units)`, total supply 100, height 0, nonce(A)=0, B absent. Genesis allocation bytes are produced by the canonical genesis encoder. A second fixture `S1` is the state after a successful `A -> B, amount=25, nonce=0`: A balance 75/nonce 1; B balance 25/nonce 0; total supply 100.
- A valid signature fixture is generated from the exact unsigned payload digest with RFC 6979 HMAC-SHA-256 and `TEST_ONLY_A`, then normalized to low-S. Store the unnormalized and normalized scalars to make malleability tests reproducible. ECDSA verification must use the digest exactly once.
- Record for each vector: fixture ID, initial state and chain ID, full wire hex, unsigned payload hex, signing preimage hex, digest hex, public key/address hex, `r/s` hex, TxID hex where canonical, stage/path, expected code/status, resulting accounts/nonce/supply/AppHash, producer version, independent verifier version, and executed command/output. Until generated, these values are **TBD**, not inferred.

The vector producer must use the fixed byte layout, secp256k1 parameters and domain strings in [transaction V1](B0_CLASSICAL_TRANSACTION_PROTOCOL_V1.md), not a library's Ethereum transaction codec. Baseline results must be checked against [state V1](B0_STATE_TRANSITION_SPEC_V1.md) and [ABCI validity V1](B0_ABCI_VALIDITY_SPEC_V1.md). Test code should compare raw bytes and full 32-byte hashes, not display abbreviations.

## 2. Required vector matrix

`CheckTx` uses the stated committed fixture; `ProcessProposal` uses its ordered block candidate; `FinalizeBlock` expectations describe the defensive decided-block path even where a correct `ProcessProposal` would reject. `State delta = none` means balances, nonces, account membership and supply do not change; finalizing a block still advances height and changes AppHash. Codes and precedence are those in the ABCI specification.

| ID | Construct / starting fixture | Expected stage and result | State and invariant outcome |
|---|---|---|---|
| V01 valid transfer | On S0, sign A→B amount 25 nonce 0. | Decode/signature/state pass; `CheckTx` Code 0, one-tx proposal ACCEPT, `FinalizeBlock` Code 0. | S1: A 75/1, B 25/0, supply 100; pending until Commit. |
| V02 incorrect signature scalar | V01 wire but replace `r` with 32 zero bytes. | Signature stage Code 4; proposal REJECT; final result Code 4. | No tx state delta. |
| V03 high-S malleability | From V01 valid low-S signature, replace `s` with `n-s`. | Signature canonicality Code 4 even if ECDSA mathematics accepts it; proposal REJECT. | No tx state delta. |
| V04 modified signed amount | Keep V01 signature but change amount 25→26 in wire. | Signature verification Code 4; proposal REJECT. | No tx state delta. |
| V05 wrong chain | Sign a canonical transfer whose wire chain ID is `pluto-other-1`, evaluate on S0 chain. | Chain comparison Code 2 before signature verification; proposal REJECT. | No tx state delta; cross-chain replay blocked. |
| V06 invalid public key | Set V01 key prefix to `0x04` while retaining 33-byte field length. | Key stage Code 3; proposal REJECT. | No tx state delta. |
| V07 malformed/noncanonical | Separate subcases: bad magic, unknown version, chain length 0/51, non-ASCII chain byte, truncated field, trailing byte, uncompressed 65-byte key changing wire length. | Parse stage Code 1 for each; proposal REJECT. | No tx state delta; no protocol TxID for noncanonical wire. |
| V08 incorrect/future nonce | Valid signature for A→B amount 25 nonce 1 on S0. | State nonce Code 6; proposal REJECT. | No tx state delta. |
| V09 replay | Apply V01, then evaluate identical V01 on S1. | State nonce Code 6; proposal REJECT. | S1 unchanged; no double spend. |
| V10 insufficient balance | On S0, sign A→B amount 101 nonce 0. | State balance Code 7; proposal REJECT. | No tx state delta; supply remains 100. |
| V11 zero amount | On S0, sign A→B amount 0 nonce 0. | Amount Code 5; proposal REJECT. | No tx state delta. |
| V12 exact-balance transfer | On S0, sign A→B amount 100 nonce 0. | Code 0; proposal ACCEPT. | A retained at balance 0/nonce 1, B balance 100/nonce 0; supply 100. |
| V13 duplicate in one block | Candidate `[V01,V01]` on S0. | `ProcessProposal` REJECT because second tx is Code 6. If anomalously finalized: results `[0,6]`. | Only first applies; S1 accounts, supply 100, block height advances. |
| V14 multiple ordered txs | On S0, `[A→B 25 nonce 0, A→B 30 nonce 1]`, both signed. | Full proposal ACCEPT, results `[0,0]`. Before either commits, committed-state-only `CheckTx` for nonce 1 returns Code 6. | Final A 45/nonce 2, B 55/nonce 0, supply 100. |
| V15 invalid proposer injection | Inject malformed V07 before canonical V01, without local `CheckTx`. | `ProcessProposal` REJECT. If anomalously finalized: `[1,0]`; first skipped. | V01 alone applies; supply 100; demonstrates independent proposal/final checks. |
| V16 sender/key binding | Replace V01 key with `P_B` but retain A signature; separately sign with B over B key while B is unfunded in S0. | First case Code 4; second case Code 7 after deriving sender B. | No ability to debit A with B key; no tx state delta. |
| V17 nonce upper bound | Synthetic, integrity-valid state fixture with A nonce `MaxUint64` and sufficient balance; sign with wire nonce `MaxUint64`. | State nonce Code 6 (cannot increment). | No wrap or tx state delta. This fixture is a boundary test, not reachable in practical block counts. |
| V18 genesis supply overflow | Canonical-looking genesis with two positive allocations whose checked sum exceeds `MaxUint64`. | `InitChain` rejects before creating state. | No V2 genesis record or AppHash accepted. |
| V19 account cap | Separate synthetic integrity-valid state with total supply at least 65,536, exactly 65,536 accounts, funded A and an absent receiver. | State account-limit Code 9; proposal REJECT. | No new account or tx state delta. |
| V20 block bounds | Separate chain fixtures: 6,145 minimum-length (146-byte) canonical records total 897,170 bytes, isolating count; 5,378 maximum-length (195-byte) records total 1,048,710 bytes, isolating byte limit. Give each transaction a valid sequential nonce or distinct funded sender as needed. | `ProcessProposal` REJECT for either bound; anomalous `FinalizeBlock` returns Code 10 once per tx. | No tx state changes; height advances on finalization. |
| V21 sender nonce after failed tx | Candidate `[insufficient-balance A nonce 0, valid A nonce 0]` on S0. | Proposal REJECT for first. Anomalous final results `[7,0]`. | Failed tx consumes no nonce; second produces S1. |
| V22 self-transfer | On S0, sign A→A amount 25 nonce 0. | Code 0; proposal ACCEPT. | A balance 100/nonce 1; one account; supply 100. |
| V23 state corruption | Flip one persisted logical account byte or V2 digest byte, then attempt reload. | State loader rejects before `Info`/proposal processing. | No corrupted state accepted or silent reset. |

A recipient-credit `uint64` overflow is mathematically unreachable from an integrity-valid, supply-bounded state when a distinct funded sender transfers without minting. The implementation must still check it defensively. A test that manufactures an inconsistent state should expect **state load/integrity rejection**, not falsely claim a reachable Code 8 transaction vector. Similarly, a height-overflow test uses a synthetic integrity-valid height `MaxInt64` and expects a block-level error without wrap.

## 3. Independent vector-generation gates

1. Freeze reviewed protocol text and record a source commit, toolchain, library versions and exact generation commands. Add any secp256k1 dependency only when first imported in authorized Phase 03B, with a separate dependency decision.
2. Generate valid bytes with the fixed fixture and a reviewed secp256k1 implementation; independently parse fields and recompute SHA-256 preimages/digests/address/TxID with another implementation or small independent calculation. Check SEC 1 point and ECDSA rules, low-S and RFC 6979 deterministic signing against primary references.
3. For every negative vector, derive it from an explicit valid or synthetic fixture, write exact mutation and error precedence, then execute `CheckTx`/proposal/finalize stages once those paths exist. Never label a state-invalid vector `CheckTx`-accepted merely because a prior transaction in the same proposed block could fund it.
4. Compare replay on two independent state instances; in Phase 03E compare four validators' ordered result codes and AppHashes. Record any skipped or blocked vector as NOT RUN with reason.
5. Keep only these public test scalars in vector artifacts. Never place generated validator/node private keys, real wallet keys or credentials in the repository.

## 4. Review questions

Which independent secp256k1 implementation will verify the Phase 03B golden vectors? Should the protocol require RFC 6979 signing for reference vectors only (proposed) or for all client signers (not remotely enforceable)? What deployment size should be used for account-cap and block-bound tests? These questions require review; no vector has been executed here.
