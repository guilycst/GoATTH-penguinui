# CLAUDE.md - GoATTH PenguinUI

## Project Overview

**GoATTH** (Go + Alpine.js + Tailwind CSS + HTMX + Templ) is a UI component library that replicates [PenguinUI](https://penguinui.com) components using Go's templating system. Hard fork of PenguinUI targeting 99.99% visual parity.

## Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.26+ | Backend server and templating |
| templ | v0.3.x | HTML template generation (.templ → .go) |
| Tailwind CSS | v4 | Utility-first styling |
| HTMX | v2.0.8 | Dynamic content loading |
| Alpine.js | v3.x | Reactive UI components |
| Playwright | v0.5700.1 | E2E testing |

## Quick Commands

```bash
# Install all dependencies
make install

# Generate templ files (REQUIRED after editing .templ files)
make generate

# Build Tailwind CSS (REQUIRED after editing CSS)
make css

# Run dev server (default port 7070)
make dev
# or with live reload:
make dev-air

# Run tests
make test          # unit tests
make test-e2e      # E2E tests (Playwright)

# Format & lint
make fmt
make lint
```

## Repository Structure

```
GoATTH-penguinui/
├── cmd/server/main.go          # Server entry point
├── components/                 # Reusable UI components
│   └── <name>/
│       ├── <name>.templ        # Component template
│       ├── types.go            # Config types and variant classes
│       └── <name>_templ.go     # Generated (DO NOT EDIT)
├── internal/
│   ├── server/server.go        # Route handlers
│   └── pages/demo/
│       ├── layout.templ        # Main layout, sidebar, theme selector
│       ├── split_view.templ    # Side-by-side comparison view
│       ├── tab_view.templ      # Tabbed comparison view
│       └── components/         # Demo pages per component
├── css/main.css                # Tailwind source + theme imports
├── all-themes.css              # 13 theme definitions
├── assets/
│   ├── js/darkmode.js          # Alpine.js dark mode store
│   └── styles.css              # Generated CSS (DO NOT EDIT)
├── tests/e2e/                  # Playwright E2E tests
└── Makefile                    # Build commands
```

## Component Development Workflow

1. **Analyze reference** — Read PenguinUI HTML in `/<component-name>/` directory
2. **Create component** — `components/<name>/types.go` + `<name>.templ`
3. **Create demo page** — `internal/pages/demo/components/<name>.templ`
4. **Register route** — Add handler in `internal/server/server.go`
5. **Add to sidebar** — Update `internal/pages/demo/layout.templ`
6. **Write E2E tests** — `tests/e2e/<name>_test.go`
7. **Build & verify** — `make generate && make dev`

## Critical Rules

### After modifying `.templ` files, ALWAYS run:
```bash
make generate
```

### After modifying CSS files, ALWAYS run:
```bash
make css
```

### Files marked "Generated" (`*_templ.go`, `assets/styles.css`) — NEVER edit manually

### Dark mode pattern — Always use both light and dark variants:
```css
bg-surface text-on-surface dark:bg-surface-dark dark:text-on-surface-dark
```

### Alpine.js arrays — NEVER initialize as `null`, always use `[]`:
```go
if string(selectedJSON) == "null" {
    selectedJSON = []byte("[]")
}
```

## Theme System

- Themes defined in `all-themes.css` using `[data-theme="name"]` selectors
- Dark mode uses `.dark` class on `<html>` via Alpine.js store
- CSS custom variant: `@custom-variant dark (&:where(.dark, .dark *))`
- Default theme: **Minimal** (black/white, no border radius)

## Component Patterns

### types.go structure:
```go
package componentname

type Config struct {
    ID    string
    Label string
    // ...
}

func (cfg Config) Classes() string { ... }
```

### .templ structure:
```templ
templ ComponentName(cfg Config) {
    // Main entry point, delegates to private helpers
}

templ helperTemplate(cfg Config) {
    // Private implementation
}
```

### Alpine.js data generation:
```go
func alpineData(cfg Config) string {
    optsJSON, _ := json.Marshal(cfg.Options)
    return fmt.Sprintf(`{ options: %s }`, optsJSON)
}
```

## Testing

- E2E tests use Playwright Go bindings
- Wait for Alpine.js: use `page.WaitForTimeout(800)` or `WaitForSelector`
- Each test should reload page for clean state
- Test visual parity between Original and GoATTH columns
- Run specific test: `make test-e2e-one TEST=TestButtonVariants`

## Current Status

**Completed (8):** Button, Accordion, Sidebar, Avatar, Badge, Banner, Card, Combobox
**Pending (16+):** Modal, Alert, Input, Select, Checkbox, Radio, Toggle, Toast, Table, Tabs, Pagination, Progress, Skeleton, Spinner, Tooltip, Dropdown, and more.
