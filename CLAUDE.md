# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CMDB (Configuration Management Database) system for enterprise IT asset management. Built with Go 1.22+ (Gin) backend and React 19 (Vite) frontend. Uses PostgreSQL with JSONB for flexible CI attribute storage and Redis for caching.

## Commands

### Backend (cmdb-api)

```bash
# Run (requires PostgreSQL and Redis)
cd cmdb-api && go run main.go

# Build
cd cmdb-api && go build -o cmdb-api

# Run tests
cd cmdb-api && go test ./...

# Run tests with coverage
cd cmdb-api && go test -cover ./...
```

### Frontend (cmdb-ui)

```bash
# Install dependencies
cd cmdb-ui && npm install

# Development server (proxies /api to localhost:8080)
cd cmdb-ui && npm run dev

# Build for production
cd cmdb-ui && npm run build

# Lint
cd cmdb-ui && npm run lint
```

### Agent (cmdb-agent)

```bash
cd cmdb-agent && go build -o cmdb-agent
```

## Architecture

### Backend: Modular Monolith

```
cmdb-api/
├── modules/
│   ├── auth/       # User authentication, JWT tokens
│   ├── core/       # CIType engine, CI instances, search
│   ├── ipam/       # IP address management
│   ├── dcim/       # Data center infrastructure
│   ├── discovery/  # Auto-discovery rules and agents
│   └── integration/# Prometheus, ELK integration
├── middleware/     # JWT auth, rate limiting
├── database/       # PostgreSQL and Redis initialization
├── pkg/            # JWT utilities, response helpers
├── config/ # Viper-based configuration
└── router/         # Gin route definitions
```

Each module follows a handler → service → repository pattern. Modules are isolated and communicate only through service interfaces.

### Frontend: React + Ant Design

```
cmdb-ui/src/
├── api/           # Axios client with JWT interceptors
├── modules/       # Page components by feature
│   ├── auth/ # Login page
│   ├── core/      # CITypeList, CIList, CIDetail
│   ├── ipam/      # SubnetTree, IPList, IPAllocate
│   ├── dcim/      # IDCList, IDCMap, RackView
│   └── discovery/ # RuleList, AgentList
├── components/    # Shared components (CIForm, RelationGraph)
├── layouts/       # AppLayout with sidebar navigation
└── router/       # React Router v6 routes
```

Frontend uses React Router v6 for routing, Ant Design 5.x for UI, Zustand for global state, and React Query for server state.

### Database Design

**CI Attribute Storage**: Uses PostgreSQL JSONB column for flexible CI attribute storage:
```sql
CREATE TABLE cmdb_core.cis (
    id BIGSERIAL PRIMARY KEY,
    type_id INT NOT NULL,
    status VARCHAR(16) DEFAULT 'active',
    attr_values JSONB NOT NULL DEFAULT '{}',
    ...
);
CREATE INDEX idx_cis_attr_gin ON cmdb_core.cis USING GIN (attr_values jsonb_path_ops);
```

**Module Schemas**: Separate PostgreSQL schemas per module (cmdb_auth, cmdb_core, cmdb_ipam, cmdb_dcim, cmdb_discovery).

## API Design

Base URL: `/api/v1`

**Authentication**: JWT Bearer token in Authorization header
```json
// Response format
{ "code": 0, "message": "success", "data": {} }
```

**Key Endpoints**:
- `POST /auth/login` - Login
- `POST /auth/register` - Register
- `GET/POST/PUT/DELETE /citypes` - CIType CRUD
- `GET/POST/PUT/DELETE /ci/:id` - CI instance CRUD
- `GET /ci/s?q=_type:Server` - CI search with query syntax
- `POST/GET /ipam/subnets` - Subnet management
- `GET /dcim/idcs` - IDC list

**CI Search Syntax**:
- `_type:Server` - Filter by CIType
- `attr:value` - Attribute filter
- `attr:(v1;v2)` - IN clause
- `attr:[a TO b]` - Range query

## Configuration

Backend uses environment variables via Viper:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB` - Redis
- `JWT_SECRET`, `JWT_EXPIRE_HOURS` - JWT settings
- `SERVER_PORT` - HTTP server port (default 8080)

## Development Notes

- Frontend dev server runs on port 3000 with API proxy to backend on port 8080
- Go module name is `cmdb-api` for backend, `cmdb-agent` for agent
- CI attributes stored as JSONB - use GIN indexes for efficient querying
- Modules in `cmdb-api/modules/` are isolated - no cross-module database access

## Infrastructure & CI/CD

### Credentials

Infrastructure credentials are stored in `~/.claude/.env.claude`. Load with:

```bash
source ~/.claude/.env.claude
```

**Jenkins**：`http://172.18.68.185:8080`
**Harbor**：idc-test-harbor.neuedu.com
**K8s Namespace**：idc-test

### Deployment Pipeline

```
feature → PR → Review → merge to dev/test → deploy to idc-test (staging)
                                                   ↓
                              pass → merge to main → deploy to idc-test (final)
```

### Docker Image Naming

```bash
# Backend
idc-test-harbor.neuedu.com/cmdb/cmdb-api:staging-v[MAJOR].[FEATURE].[BUGFIX]

# Frontend
idc-test-harbor.neuedu.com/cmdb/cmdb-ui:staging-v[MAJOR].[FEATURE].[BUGFIX]
```

### Version Management

Version files (`cmdb-api/version`, `cmdb-ui/version`) store current version. Pipeline auto-increments on each deploy.

### Helm Chart

```bash
k8s/cmdb/                    # Helm chart
  ├── Chart.yaml             # Chart metadata
  ├── values.yaml            # Default values
  ├── values-staging.yaml    # Staging config (namespace: idc-test)
  └── templates/             # K8s resources
```

Example: `idc-test-harbor.neuedu.com/cmdb/cmdb-api:staging-v1.0.0`
