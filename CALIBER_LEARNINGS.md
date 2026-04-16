# Caliber Learnings

Accumulated patterns and anti-patterns from development sessions.
Auto-managed by [caliber](https://github.com/caliber-ai-org/ai-setup) — do not edit manually.

- **[gotcha:project]** For local live-reload with Docker Compose, Go binary changes need the compose watch entry to use `sync+restart` and a concrete `target` path inside the container; `rebuild` was not the right action for this workflow.
- **[correction:project]** When the user says a related component change is fine, it can be updated too — don’t assume only the originally mentioned service is in scope if the user explicitly broadens it.
- **[gotcha:project]** `go tool <name>` only works if the tool is actually listed in the module’s `tool` block and the dependency has been downloaded; if it isn’t available, run `go mod download` first and don’t rely on `go tool -n` to prove the tool is installed.
- **[gotcha:project]** The GoATTH `goatth` tool is published as a module tool, but `go tool goatth` can still fail in the repo if the current module/workspace context isn’t resolving it as expected; verify the module dependency and workspace context before assuming the tool exists.
- **[gotcha:project]** `go.work` in the repo root only includes a subset of modules; if a tool or dependency is in a different module, run commands from the correct module directory rather than assuming the workspace covers it.
- **[gotcha:project]** The GoATTH assets embed no longer includes an `images` directory; if a tool or build step expects embedded images, check `assets/embed.go` before assuming they’re part of the packaged assets.
- **[pattern:project]** When a generated or vendored helper command is missing, inspect `go.mod`’s `tool` block, then use `go mod download` and re-check the command from the owning module directory rather than trying to install it manually.
- **[correction:project]** If a user asks for a handoff, use the existing project handoff doc as the source of truth and update it instead of inventing a new summary from scratch.
- **[correction:project]** When the user asks to commit/push a dependency update, include the module tag/versioning step as part of the task; don’t stop at editing `go.mod`/`go.sum`.
- **[gotcha:project]** For kubeconfig/OIDC flows, a context that looks valid can still fail with "No valid id-token" immediately after `kubectl config set-credentials`; verify the auth-provider flow itself, not just that the context was created.
