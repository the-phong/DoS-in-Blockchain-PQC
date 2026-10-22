# B0 classical transaction protocol V1 — Phase 03A proposal

**Status: PROPOSED for controller and owner review; not implemented or accepted.** This document defines consensus bytes for a minimal, fee-free, account-based value transfer. The current Phase 02 application still rejects every transaction. Phases 03B–03E need separate authorization. Normative words **MUST**, **MUST NOT**, and **SHOULD** apply to a future B0 implementation only.

## 1. Scope and choices

| Decision | Proposed normative choice | Rationale and trade-off |
|---|---|---|
| Transaction purpose | Transfer indivisible B0 units between 20-byte accounts only. | Small correctness baseline; no fee market, gas, memo, expiry, EVM call or contract data. Spam resistance is absent and must be addressed at a later approved gate, without changing consensus validity using local resources. |
| Version | Magic `PLUTOTX1` plus unsigned big-endian `uint16` version `1`. | Explicit format boundary; unknown versions fail closed. |
| Chain binding | The signed payload carries the Phase 02-compatible chain ID: 1–50 ASCII bytes from `[A-Za-z0-9_.-]`; it MUST equal committed state/genesis chain ID byte-for-byte. | Prevents cross-chain replay, subject to operators using distinct chain IDs. |
| Sender identity | A compressed secp256k1 public key is carried; the sender account address is deterministically derived from that key. No independent sender-address field. | Avoids a second, possibly conflicting identity field. Address collisions remain a generic 160-bit hash assumption. |
| Signature | ECDSA over secp256k1, 64-byte compact `r || s`, low-S only, over a SHA-256 domain-separated digest. | Unique S representative and fixed size; this is neither Ethereum RLP/EIP-155 nor an Ethereum address/signature. |
| Nonce | Per-sender `uint64`, initially zero; a successful transaction requires equality with current account nonce and increments it by one. | Ordered replay protection; no future-nonce queue semantics in B0 admission. |
| Transaction ID | Domain-separated SHA-256 of the entire canonical signed wire record. | Binds signature bytes; differently signed forms of one intent can have different IDs, but account nonce prevents repeated execution. This is distinct from any CometBFT RPC transaction hash. |

No fee, gas price, memo, expiry, sender-address copy, key version or EVM field is necessary for the B0 value-transfer invariant. Adding any later requires a new protocol version and separate decision; a future EVM adapter must map this protocol explicitly rather than pretend Ethereum wire compatibility.

## 2. Exact wire encoding

All integers are **unsigned**, fixed-width, **big-endian**. `ASCII` denotes exact bytes, not a Unicode or JSON representation. Concatenation `||` has no padding. `N = len(chain_id)` must be 1–50. `wire = unsigned_payload || signature` and has exactly `145 + N` bytes (146–195 inclusive); any shorter, longer or trailing bytes are invalid. No variable-length integer, optional field or alternate encoding exists.

| Offset (zero-based) | Width | Field | Rule |
|---|---:|---|---|
| `0..7` | 8 | magic | ASCII `PLUTOTX1` |
| `8..9` | 2 | version | `0x0001` |
| `10..11` | 2 | chain ID length | `uint16_be(N)`; 1–50 |
| `12..11+N` | N | chain ID | exact permitted ASCII bytes |
| `12+N..44+N` | 33 | sender public key | SEC 1 compressed point, prefix `0x02` or `0x03` |
| `45+N..64+N` | 20 | receiver address | uninterpreted raw 20 bytes, including all-zero if supplied |
| `65+N..72+N` | 8 | amount | `uint64_be`, in indivisible B0 units; must be greater than zero |
| `73+N..80+N` | 8 | nonce | `uint64_be` |
| `81+N..112+N` | 32 | `r` | big-endian ECDSA scalar |
| `113+N..144+N` | 32 | `s` | big-endian ECDSA scalar |

`unsigned_payload` ends after byte `80+N` and has length `81+N`. The receiver has no checksum in wire form. A textual form, if added later, is an RPC/UI concern; hex in documents is a display of these raw bytes and is never signed as text. The wire length bound includes all fields and is checked before allocation or parsing. Re-encoding a decoded transaction MUST reproduce its input bytes exactly.

## 3. Signed intent and signature validity

Let `SIGN_DOMAIN = ASCII("PLUTO/B0/SIGN/V1") || 0x00`. Let `U` be the exact `unsigned_payload` bytes above. The 32-byte signing digest is:

```text
D = SHA-256(SIGN_DOMAIN || uint32_be(len(U)) || U)
```

The chain ID, sender key, receiver, amount and nonce are all covered by `D`. No signature byte is included. The domain and version prevent interpreting the same bytes as a different Pluto operation or future hybrid intent. The ECDSA message representative is the unsigned big-endian integer represented by the 32-byte `D`; secp256k1 ECDSA scalar arithmetic is modulo `n` as in SEC 1. The verifier supplies exactly `D` to a prehashed-message ECDSA API and MUST NOT SHA-256 it a second time. RFC 6979 reference signing, when used, takes this same 32-byte digest as its hashed-message input.

Curve parameters are secp256k1 from [SEC 2 v2.0, section 2.4.1](https://www.secg.org/sec2-v2.pdf): field prime `p = FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F` and group order `n = FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141` (hex). A 33-byte public key MUST have prefix `02` or `03`; its 32-byte x coordinate MUST be `< p`. Point decompression MUST produce a finite point on `y² = x³ + 7 (mod p)` with the prefix's y parity. The infinity, off-curve, non-canonical x, uncompressed (`04`), hybrid and trailing-byte encodings MUST be rejected. The group has cofactor 1; complete on-curve validation implies subgroup membership for secp256k1. See [SEC 1 v2.0, sections 2.3.3 and 3.2.2](https://www.secg.org/sec1-v2.pdf).

The 32-byte `r` and `s` values MUST each be in `[1,n-1]`; `s` MUST additionally satisfy `s <= floor(n/2)`. Values with leading zero bytes are valid only as part of the required fixed 32-byte scalar width; DER, shorter integers, recovery IDs and high-S alternatives are invalid. After scalar checks, verify ECDSA `(r,s)` against the validated public key and `D` using SEC 1. A signer SHOULD use [RFC 6979](https://www.rfc-editor.org/rfc/rfc6979) with HMAC-SHA-256 to avoid fragile signing randomness, then normalize `s := min(s, n-s)`. Verification cannot and does not infer the signer's nonce-generation method; any valid low-S signature is accepted.

## 4. Address and transaction ID

Let `ADDR_DOMAIN = ASCII("PLUTO/B0/ADDR/V1") || 0x00`. For a validated compressed public key `P`, the sender address is the **first 20 bytes** of `SHA-256(ADDR_DOMAIN || P)`, in digest byte order. No Keccak hash, EIP-55 checksum or Ethereum address rule is used. A receiver address is any raw 20-byte value; funds sent to an address without a known key may be unspendable, so clients SHOULD derive and verify recipients before submission.

Let `ID_DOMAIN = ASCII("PLUTO/B0/TXID/V1") || 0x00`. For canonical full `wire` bytes, `TxID = SHA-256(ID_DOMAIN || uint32_be(len(wire)) || wire)`. A protocol TxID is defined once structural decoding and scalar range/low-S canonicality pass; ECDSA verification or state validity may still fail. Malformed or noncanonical bytes have no protocol TxID. TxID is an identifier, not a second authentication check or a state commitment. Its 32 bytes are normally displayed as lowercase hex without altering the bytes.

## 5. Validation order and error precedence

The future common evaluator MUST apply these stages in order and stop at the first failure, using the codes in the [ABCI validity specification](B0_ABCI_VALIDITY_SPEC_V1.md):

1. Bound input length to 146–195, then check magic, version, chain length/characters, exact calculated total length and no trailing data. Reject unknown versions; never reinterpret them as V1.
2. Compare wire chain ID to the current committed/candidate chain ID.
3. Validate the compressed sender public key as above.
4. Parse `r,s`, enforce range and low-S, compute `D`, and verify ECDSA.
5. Require positive amount.
6. Derive sender address and compare the wire nonce to that account's current nonce.
7. Require sender balance at least amount.
8. If sender and receiver differ, require receiver balance plus amount to fit `uint64`; if receiver is absent, require room under the account-count cap.
9. Only then apply the state transition defined in [B0 state transitions](B0_STATE_TRANSITION_SPEC_V1.md).

Structural fields are decoded before signature verification, but amount is deliberately checked after cryptographic validity so each invalid record has one reproducible error precedence. Checks 1–5 are state-independent except chain matching; checks 6–8 use the sequential candidate state. An all-zero receiver address is a regular address, not a burn address. There is no fallback to another signature scheme.

## 6. Review questions and references

The owner/controller should approve or revise these proposed choices before Phase 03B: (Q1) fee-free B0 units and their human display scale; (Q2) 20-byte SHA-256-derived addresses versus a longer digest; (Q3) account/block limits and future capacity upgrade path; (Q4) whether Phase 03C should allow only fresh V2 genesis or separately specify a migration from Phase 02 homes. No choice here is marked accepted by the owner.

Primary references: [SEC 1 v2.0](https://www.secg.org/sec1-v2.pdf), [SEC 2 v2.0](https://www.secg.org/sec2-v2.pdf), [RFC 6979](https://www.rfc-editor.org/rfc/rfc6979), and the pinned [CometBFT v0.38.26 ABCI application requirements](https://github.com/cometbft/cometbft/blob/v0.38.26/spec/abci/abci%2B%2B_app_requirements.md). None of these sources makes this proposed Pluto wire format a standard.
