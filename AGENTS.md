# Repository Guidelines

## Project Structure & Module Organization
This Go 1.23.3 workspace centers on `main.go` and the `pkg/` tree, where application services (for example `pkg/server/organization_iam`) and repositories (`pkg/db`) live alongside unit tests (`*_test.go`). Proto contracts are under `proto/<domain>/`, generated code sits next to its `.proto`, and helper scripts land in `hack/` (mock generation, tooling install). Database schema changes are stored in `migration/`, while infrastructure bits (Dockerfiles, CodeBuild specs) stay at the repo root for easy CI wiring.

## Build, Test, and Development Commands
- `make install` — installs grpc/protoc plugins expected by the proto targets.
- `make proto` — formats proto files via `clang-format`, then runs both validated and non-validated `protoc` passes plus mock generation; ensure `protoc`, `protoc-gen-go`, and `protoc-gen-validate` are on `PATH`.
- `GO111MODULE=on go test ./...` (or `make test`) — executes the full Go test suite.
- `make lint` — runs `golangci-lint` with a 5‑minute timeout; keep GOFLAGS free of `-mod=vendor`.
- `sh hack/generate-mock.sh <path>` — regenerates gomock stubs for a package (e.g., `proto/organization_iam` or `pkg/db`).

## Coding Style & Naming Conventions
Use standard `gofmt`/`goimports` formatting (tabs, camelCase identifiers, exported names with package prefixes). Proto definitions follow snake_case field names with explicit `go_package` options; run `make fmt` before committing proto edits. Prefer descriptive filenames (e.g., `access_token.go`) and keep package names singular. When editing generated code, always update the source `.proto` instead of editing `.pb.go` files directly.

## Testing Guidelines
Unit tests live beside the code they cover using the `testing` package plus `stretchr/testify`, `go-sqlmock`, and helpers in `pkg/test`. Name tests `Test<Feature>` or `Test<Struct>_<Behavior>` for clarity. Run `go test ./pkg/...` to scope to server logic, or `go test ./proto/...` after regenerating validators. Aim to exercise happy paths, gRPC validation failures, and database error branches; mock repositories via `pkg/db/mocks`.

## Commit & Pull Request Guidelines
Recent history (`feat: …`, `fix: …`) shows a Conventional Commits style—start messages with a lowercase type and keep the subject imperative within ~72 chars. Pull requests should describe motivation, link Jira/GitHub issues, and include verification steps (tests, lint, proto regeneration). When UI or API behavior changes, attach gRPC or HTTP examples plus any migration considerations. Ensure each PR leaves the tree buildable (`make proto && make test`) before requesting review.

## Security & Configuration Tips
Store secrets in environment variables or parameter stores; avoid committing credentials under `config/` or scripts. When tweaking IAM or token logic, review related protobufs and migrations to ensure hashes, expirations, and role bindings remain consistent across `pkg/db` and `pkg/server`.
