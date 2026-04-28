# Repository Guidelines

## Project Structure & Module Organization
This Go 1.23.3 workspace centers on `main.go` and the `pkg/` tree. Unit tests (`*_test.go`) live beside the code they cover. Database schema changes are stored in `migration/`, while infrastructure bits (Dockerfiles, CodeBuild specs) stay at the repo root for easy CI wiring.

### Directory Structure
```
pkg/
├── server/     # gRPC service implementations
├── db/         # Database repositories
├── model/      # GORM data models
├── ai/         # AI service integration
└── test/       # Test utilities and mocks

proto/          # Protocol buffer definitions (generated code sits next to its .proto)
hack/           # Build scripts (mock generation, tooling install)
```

### Consistent Naming Convention Across Packages
Files and directories in `proto/`, `pkg/model/`, `pkg/server/`, and `pkg/db/` must all use the same name for a given service.
- `proto/organization/` → `pkg/model/organization.go`, `pkg/server/organization/`, `pkg/db/organization.go`
- `proto/alert/` → `pkg/model/alert.go`, `pkg/server/alert/`, `pkg/db/alert.go`
- `proto/org_alert/` → `pkg/model/org_alert.go`, `pkg/server/org_alert/`, `pkg/db/org_alert.go`
- `proto/finding/` → `pkg/model/finding.go`, `pkg/server/finding/`, `pkg/db/finding.go`

## Build, Test, and Development Commands
- `make install` — installs repository tooling including grpc/protoc-related dependencies expected by the build targets.
- `make generate-mock` — regenerates mocks across `proto/`, `pkg/db`, and `pkg/ai`; use `sh hack/generate-mock.sh <path>` for a narrower target.
- `GO111MODULE=on go test ./...` (or `make test`) — executes the full Go test suite.
- `make lint` — runs `golangci-lint` with a 5‑minute timeout; keep GOFLAGS free of `-mod=vendor`.
- `make build` / `make build-ci` — builds the Docker image locally or for CI; `make build` runs tests first.

## Coding Style & Naming Conventions
Use standard `gofmt`/`goimports` formatting (tabs, camelCase identifiers, exported names with package prefixes). Prefer descriptive filenames (e.g., `access_token.go`) and keep package names singular.

## Testing Guidelines
Unit tests live beside the code they cover using the `testing` package plus `stretchr/testify`, `go-sqlmock`, and helpers in `pkg/test`. Name tests `Test<Feature>` or `Test<Struct>_<Behavior>` for clarity. Run `go test ./pkg/...` to scope to server logic. Aim to exercise happy paths, gRPC validation failures, and database error branches; mock repositories via `pkg/db/mocks`.

## Commit & Pull Request Guidelines
Recent history (`feat: …`, `fix: …`) shows a Conventional Commits style—start messages with a lowercase type and keep the subject imperative within ~72 chars. Pull requests should describe motivation, link Jira/GitHub issues, and include verification steps (tests, lint, proto regeneration). When UI or API behavior changes, attach gRPC or HTTP examples plus any migration considerations. Ensure each PR leaves the tree buildable (`make proto && make test`) before requesting review.

## Security & Configuration Tips
Store secrets in environment variables or parameter stores; avoid committing credentials under `config/` or scripts. When tweaking IAM or token logic, review related protobufs and migrations to ensure hashes, expirations, and role bindings remain consistent across `pkg/db` and `pkg/server`. AI-related settings are also environment-driven, including `OPENAI_TOKEN`, `CHATGPT_MODEL`, and notification credentials such as `SLACK_API_TOKEN`.

### Authorization
Each gRPC handler in Core is authorized by `ca-risken/gateway` using `organization_id` and `project_id`. Core itself does not implement authorization logic; the Gateway is responsible for request authorization.

## Translation Guidelines
- Organizationの日本語訳は"組織"とせず、"Org"にしてください。

## CLAUDE.md Placement Rule
Every directory that contains an `AGENTS.md` file must also have a `CLAUDE.md` file that includes `@AGENTS.md` to ensure the agent guidelines are loaded automatically.
