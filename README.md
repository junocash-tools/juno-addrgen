# juno-addrgen

Offline address derivation (UFVK + index -> `j1...`) for Juno Cash.

## Usage

- Derive one address:
  - `juno-addrgen derive --ufvk <jview1...> --index 0`
- Derive a batch:
  - `juno-addrgen batch --ufvk <jview1...> --start 0 --count 10`
- JSON output:
  - add `--json`

UFVKs are sensitive (watch-only, but reveal incoming transaction details). Avoid logging or sharing them.

## Build & test

Requirements: Go + Rust toolchain.

- Build: `make build`
- Test (unit + integration + e2e): `make test`
