# Proto Guidelines

## Scope
`proto/` contains protocol buffer definitions, and generated code sits next to each `.proto` file.

## Editing Rules
Proto definitions follow snake_case field names with explicit `go_package` options.
Run `make fmt` before committing proto edits.
When editing generated code, always update the source `.proto` instead of editing `.pb.go` files directly.

## Generation
- `make proto` formats proto files via `clang-format`, then runs both validated and non-validated `protoc` passes plus mock generation.
- `make proto-without-validate` / `make proto-validate` run each protobuf generation pass separately when debugging generation failures.

This repository intentionally uses the legacy `--go_out=plugins=grpc,paths=source_relative:proto` flow.
Generated service and message code stays in a single `*.pb.go` file; do not expect separate `*_grpc.pb.go` outputs.
If generation starts differing across environments, ensure the required repository tooling is installed.

## Testing
Run `go test ./proto/...` after regenerating validators.
