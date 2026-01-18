# juno-addrgen

Offline address derivation (UFVK + index -> `j1...`) for Juno Cash.

## API stability

- For automation/integrations, treat `--json` output as the stable API surface. Human-oriented output may change.
- JSON outputs are versioned (see below) to allow additive evolution without breaking consumers.

## Usage

- Derive one address:
  - `juno-addrgen derive --ufvk <jview1...> --index 0`
- Derive a batch:
  - `juno-addrgen batch --ufvk <jview1...> --start 0 --count 10`
- JSON output:
  - add `--json`
 - Read UFVK from a file:
   - `juno-addrgen derive --ufvk-file ./ufvk.txt --index 0`
 - Read UFVK from an env var name:
   - `juno-addrgen derive --ufvk-env JUNO_UFVK --index 0`

UFVKs are sensitive (watch-only, but reveal incoming transaction details). Avoid logging or sharing them.

Notes:

- `--uvfk` is accepted as an alias for `--ufvk`.
- Exchanges must persistently map derived deposit addresses (or their derivation indices) to internal accounts; on-chain data is encrypted and you cannot “match addresses” the Bitcoin way.

## JSON output

All JSON responses include:

- `version`: response schema version (string, currently `"v1"`)
- `status`: `"ok"` or `"err"`

Derive (`derive --json`):

```json
{ "version": "v1", "status": "ok", "address": "j1..." }
```

Batch (`batch --json`):

```json
{ "version": "v1", "status": "ok", "start": 0, "count": 10, "addresses": ["j1...", "..."] }
```

Errors:

```json
{ "version": "v1", "status": "err", "error": "ufvk_invalid_bech32m", "message": "..." }
```

## Build & test

Requirements: Go + Rust toolchain.

- Build: `make build`
- Test (unit + integration + e2e): `make test`
