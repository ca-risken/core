# Proto Guidelines

## Scope
`proto/` contains protocol buffer definitions, and generated code sits next to each `.proto` file.

## Naming
Files and directories in `proto/`, `pkg/model/`, `pkg/server/`, and `pkg/db/` must use the same name for a given service.
- `proto/organization/` -> `pkg/model/organization.go`, `pkg/server/organization/`, `pkg/db/organization.go`
- `proto/alert/` -> `pkg/model/alert.go`, `pkg/server/alert/`, `pkg/db/alert.go`
- `proto/org_alert/` -> `pkg/model/org_alert.go`, `pkg/server/org_alert/`, `pkg/db/org_alert.go`
- `proto/finding/` -> `pkg/model/finding.go`, `pkg/server/finding/`, `pkg/db/finding.go`

## Editing Rules
Proto definitions follow snake_case field names with explicit `go_package` options.
Run `make fmt` before committing proto edits.
When editing generated code, always update the source `.proto` instead of editing `.pb.go` files directly.

## Generation
- `make install` installs grpc/protoc plugins expected by the proto targets.
- `make proto` formats proto files via `clang-format`, then runs both validated and non-validated `protoc` passes plus mock generation.
- `make proto-without-validate` / `make proto-validate` run each protobuf generation pass separately when debugging generation failures.
- `make generate-mock` regenerates all mocks; use `sh hack/generate-mock.sh <path>` for a narrower target.

This repository intentionally uses the legacy `--go_out=plugins=grpc,paths=source_relative:proto` flow.
Generated service and message code stays in a single `*.pb.go` file; do not expect separate `*_grpc.pb.go` outputs.
If generation starts differing across environments, run `make install` first to align the expected plugin setup.

## Testing
Run `go test ./proto/...` after regenerating validators.
