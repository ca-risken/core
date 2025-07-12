# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Test
```bash
# Run tests
make test

# Lint code
make lint

# Build Docker image (runs tests first)
make build

# Install dependencies
make install
```

### Code Generation
```bash
# Generate all protobuf files and mocks
make proto

# Generate mocks only
make generate-mock

# Format proto files
make fmt
```

**Protocol Buffer Generation**:
- Uses `plugins=grpc` approach for unified file generation
- Compatible with protoc-gen-go v1.33.0 while maintaining old format
- No separate `*_grpc.pb.go` files - all code in `*.pb.go`
- mockery v2.53.3 for consistent mock generation
- Run `make install` to install required tools with correct versions

### Individual Services
```bash
# Generate proto without validation
make proto-without-validate

# Generate proto with validation
make proto-validate

# Generate specific mocks
make proto-mock      # Service mocks
make repository-mock # Database mocks
make ai-mock        # AI service mocks
```

### Docker Operations
```bash
# Build for CI
make build-ci

# Push/pull images
make push-image
make pull-image

# Multi-arch manifest
make create-manifest
make push-manifest
```

## Architecture Overview

### Service Structure
RISKEN Core is a monolithic Go application hosting multiple gRPC services:

- **Finding Service**: Security findings, resources, AI summaries
- **Alert Service**: Alert conditions, notifications, analysis
- **IAM Service**: Authentication, authorization, user management
- **AI Service**: OpenAI integration for security analysis and report generation
- **Project Service**: Project management and organization
- **Report Service**: Security reporting and data aggregation
- **Organization Services**: Multi-tenant organization management

### Key Architectural Patterns

**Service Communication**: 
- All services run in single process but communicate via gRPC clients
- Protocol Buffers define service contracts in `proto/`
- Services can act as both client and server

**Database Architecture**:
- MySQL with master/slave configuration
- GORM for ORM with repository pattern
- Connection pooling with configurable limits

**AI Integration**:
- OpenAI SDK with streaming and non-streaming responses
- FreeCache for response caching (200MB, 1-hour TTL)
- Configurable model selection (default: gpt-4o-mini)
- Multilingual support with specialized security prompts
- SQL report generation tool with security constraints

**Authorization**:
- Multi-level IAM: User → Role → Policy → Resource
- Project-based resource isolation
- Access token authentication
- Policy regex pattern matching

### Directory Structure
```
pkg/
├── server/     # gRPC service implementations
├── db/         # Database repositories
├── model/      # GORM data models
├── ai/         # AI service integration
└── test/       # Test utilities and mocks

proto/          # Protocol buffer definitions
hack/           # Build scripts
```

### Development Guidelines

**Adding New Services**:
1. Define service contract in `proto/`
2. Run `make proto` to generate Go code
3. Implement service in `pkg/server/`
4. Add repository layer in `pkg/db/`
5. Create models in `pkg/model/`
6. Generate mocks with `make generate-mock`

**Testing Requirements**:
- All services must have unit tests
- Use generated mocks for external dependencies
- Database operations use go-sqlmock
- HTTP operations use jarcoal/httpmock
- Test utilities available in `pkg/test/`

**Performance Considerations**:
- Use batch operations for bulk data (e.g., `PutFindingBatch`)
- Implement caching for expensive operations
- Use streaming responses for large datasets
- Consider connection pooling for external services

**Security Requirements**:
- All endpoints must implement authorization checks
- Use IAM service for permission validation
- Project-based resource scoping required
- Admin users can bypass project restrictions

### Configuration

Environment variables control service behavior:
- `OPENAI_TOKEN`: Required for AI features
- `CHATGPT_MODEL`: AI model selection
- `SLACK_API_TOKEN`: Notification integration
- `DB_*`: Database connection parameters
- `PROFILE_*`: Profiling and tracing settings

### Common Development Tasks

**Adding New Finding Types**:
- Update finding model in `pkg/model/`
- Add repository methods in `pkg/db/`
- Update AI prompts in `pkg/ai/finding.go`
- Add alert rule support in alert service

**Extending AI Capabilities**:
- Modify prompts in `pkg/ai/finding.go`
- Add new methods to AIService interface
- Update caching keys for new operations
- Consider streaming vs non-streaming responses
- Add SQL generation tools for report functionality

**Adding New Alert Types**:
- Update alert condition model
- Add matching logic in alert service
- Update notification templates
- Add new notification channels as needed

**AI Report Generation**:
- GenerateReport API accepts natural language prompts
- Generates SQL queries against Finding tables
- Security constraints: project_id filtering, SELECT-only queries
- Results limited to 1000 records for performance
- Excludes PendFinding table for security

### Service Dependencies

**Alert Service Dependencies**:
- Finding Service (for alert analysis)
- Project Service (for project validation)
- IAM Service (for notification authorization)

**Finding Service Dependencies**:
- AI Service (for summaries)
- None for core finding operations

**IAM Service Dependencies**:
- Finding Service (for cleanup operations)
- Self-contained for authorization

**AI Service Dependencies**:
- Finding Repository (for report generation)
- OpenAI API (for language model integration)
- Database connection (for SQL execution)

This architecture supports scaling by adding new services while maintaining clear service boundaries and consistent patterns.
