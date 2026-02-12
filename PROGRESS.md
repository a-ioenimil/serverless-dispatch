# Progress Log

Last updated: 2026-02-12

## Scope
- Role-based task dispatching for Admin and Member users.
- Frontend Kanban board (TanStack Router + React Query) with AWS API Gateway backend.
- Cognito auth flow with protected routes.
- Email notifications for task assignment and status lifecycle updates.

## Implemented
- Admin-only task assignment support in board create/edit interactions. See [frontend/src/components/kanban/TaskBoard.tsx](frontend/src/components/kanban/TaskBoard.tsx).
- Username added to signup flow and auth user model for assignment identity. See [frontend/src/components/mvpblocks/signup-form-3.tsx](frontend/src/components/mvpblocks/signup-form-3.tsx) and [frontend/src/lib/auth.ts](frontend/src/lib/auth.ts).
- Kanban board with three status columns (OPEN, IN_PROGRESS, DONE), drag-and-drop, and optimistic updates. See [frontend/src/components/kanban/TaskBoard.tsx](frontend/src/components/kanban/TaskBoard.tsx).
- Existing Cognito auth integration and protected route gating retained. See [frontend/src/integrations/auth/protected-route.tsx](frontend/src/integrations/auth/protected-route.tsx) and [frontend/src/integrations/auth/auth-context.tsx](frontend/src/integrations/auth/auth-context.tsx).

## Backend Functions
- API Gateway routes (HTTP API):
  - POST /tasks -> create task
  - GET /tasks -> list tasks
  - PUT /tasks/{taskId} -> update task
  See [terraform/infra/modules/api_gateway/main.tf](terraform/infra/modules/api_gateway/main.tf).
- Lambda handlers:
  - Create task: [functions/cmd/api-task-create/main.go](functions/cmd/api-task-create/main.go)
  - List tasks: [functions/cmd/api-task-list/main.go](functions/cmd/api-task-list/main.go)
  - Update task: [functions/cmd/api-task-update/main.go](functions/cmd/api-task-update/main.go)
  - Async notifier: [functions/cmd/async-notifier/main.go](functions/cmd/async-notifier/main.go)
  - Auth post-signup trigger: [functions/cmd/auth-post-signup/main.go](functions/cmd/auth-post-signup/main.go)
- Domain and services:
  - Task domain model and statuses: [functions/internals/task/domain/task.go](functions/internals/task/domain/task.go)
  - DTOs and request types: [functions/internals/task/services/DTOs.go](functions/internals/task/services/DTOs.go)
  - Task service RBAC logic: [functions/internals/task/services](functions/internals/task/services)
  - Notification service supports status change + reassignment alerts: [functions/internals/notification/services/notifier.go](functions/internals/notification/services/notifier.go)
  - Recipient resolver for admin group/user email lookup: [functions/internals/notification/ports/recipient_resolver.go](functions/internals/notification/ports/recipient_resolver.go) and [functions/internals/identity/infrastructure/cognito/user_directory.go](functions/internals/identity/infrastructure/cognito/user_directory.go)
- Storage:
  - DynamoDB repository implementation: [functions/internals/task/infrastructure/dynamodb](functions/internals/task/infrastructure/dynamodb)

## Infra & CI/CD
- Async notifier now receives `USER_POOL_ID` via Terraform compute module env vars. See [terraform/infra/modules/compute/functions_notify.tf](terraform/infra/modules/compute/functions_notify.tf).
- Root infra and compute modules now include `user_pool_id` variable wiring. See [terraform/infra/main.tf](terraform/infra/main.tf), [terraform/infra/variables.tf](terraform/infra/variables.tf), and [terraform/infra/modules/compute/variables.tf](terraform/infra/modules/compute/variables.tf).
- CI/CD workflow now injects `TF_VAR_user_pool_id` from `COGNITO_USER_POOL_ID` secret and fails fast if missing. See [.github/workflows/deploy.yml](.github/workflows/deploy.yml).
- Required CI/CD secrets documented in [README.md](README.md).

## Configuration
- Frontend env vars:
  - VITE_API_BASE_URL (API Gateway base URL)
  - VITE_COGNITO_CLIENT_ID
  - VITE_COGNITO_USER_POOL_ID
  - VITE_AWS_REGION
  Example is in [frontend/.env.example](frontend/.env.example).
- Backend notifier env vars:
  - FROM_EMAIL
  - USER_POOL_ID
- Terraform vars:
  - user_pool_id (for async notifier identity lookups)

## How to Run
- Frontend:
  - npm install
  - npm run dev
  See [frontend/README.md](frontend/README.md).

## Validation Status
- Backend: `go test ./...` passes in [functions](functions).
- Terraform: `terraform validate` passes in [terraform/infra](terraform/infra) (with non-blocking deprecation warnings from module/provider internals).
- Frontend: build still reports pre-existing TypeScript issues outside this scope plus one drag typing issue in board rendering; requires a separate cleanup pass before full green build.

