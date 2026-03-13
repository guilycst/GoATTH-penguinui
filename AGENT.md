# AGENT.md - GoATTH PenguinUI

## Project Overview

**GoATTH** (Go + Alpine.js + Tailwind CSS + HTMX + Templ) is a UI component library that replicates [PenguinUI](https://penguinui.com) components using Go's templating system. This is a hard fork of PenguinUI with the goal of achieving 99.99% visual parity.

### Key Features
- Side-by-side comparison view: Original HTML vs GoATTH components
- Tab view layout with GoATTH first, Original second
- Complete theme system with 13 built-in themes (Minimal, Arctic, Modern, etc.)
- Global dark mode toggle with Alpine.js store
- HTMX integration for dynamic fragment reloading
- Comprehensive E2E tests with Playwright

## Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.26+ | Backend server and templating |
| templ | v0.3.x | HTML template generation |
| Tailwind CSS | v4 | Utility-first styling |
| HTMX | v2.0.8 | Dynamic content loading |
| Alpine.js | v3.x | Reactive UI components |
| Playwright | - | E2E testing |

## Repository Structure

```
GoATTH-penguinui/
├── assets/
│   ├── js/
│   │   └── darkmode.js          # Alpine.js dark mode store
│   └── styles.css               # Generated Tailwind CSS
├── components/
│   └── button/
│       ├── button.templ         # Button component template
│       └── types.go             # Button types and variant classes
├── css/
│   └── main.css                 # Tailwind source + theme definitions
├── internal/
│   └── pages/
│       └── demo/
│           ├── components/      # Demo components
│           ├── layout.templ     # Main layout with theme selector
│           ├── split_view.templ # Side-by-side comparison view
│           └── tab_view.templ   # Tabbed comparison view
├── all-themes.css               # Complete theme definitions (13 themes)
├── go.mod                       # Go module definition
├── Makefile                     # Build commands
└── justfile                     # Just command runner (optional)
```

## Build Commands

```bash
# Generate templ files (compile .templ to .go)
make generate

# Build CSS (process Tailwind)
make css

# Run development server
make dev                    # Or: go run cmd/server/main.go -port 7070

# Run all tests
make test

# Run E2E tests
make test-e2e              # Or: cd tests/e2e && go test -v
```

## Development Workflow

### Starting Development

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Generate templ files:**
   ```bash
   make generate
   ```

3. **Build CSS:**
   ```bash
   make css
   ```

4. **Run server:**
   ```bash
   go run cmd/server/main.go -port 7070
   ```

5. **Open browser:**
   http://localhost:7070/components/button

### Making Changes

**After modifying .templ files:**
```bash
make generate
```

**After modifying CSS:**
```bash
make css
```

**Both:**
```bash
make generate && make css
```

## Theme System Architecture

### How It Works

GoATTH uses PenguinUI's theme system exactly:

1. **Theme Variables:** Defined in `all-themes.css` within `@layer base` using `[data-theme="name"]` selectors
2. **Color Pattern:** Each theme defines light colors (e.g., `--color-surface`) and dark variants (e.g., `--color-surface-dark`)
3. **Dark Mode:** Uses `@custom-variant dark` with `.dark` class on `<html>` element
4. **Component Classes:** Use `dark:` prefixes that automatically switch when `.dark` class is present

### Theme Files

- **`css/main.css`**: Source file with `@import "tailwindcss"` and theme structure
- **`all-themes.css`**: Complete theme definitions (Minimal, Arctic, Modern, High-Contrast, etc.)
- **`assets/js/darkmode.js`**: Alpine.js store managing dark mode state

### Default Theme

The default theme is **Minimal** (black/white with no border radius).

### Dark Mode Toggle

The dark mode toggle is in `internal/pages/demo/layout.templ` and uses Alpine.js store:

```javascript
// Click handler: @click="$store.darkMode.toggle()"
// State check: x-show="$store.darkMode.on"
```

## Key Files Reference

### Component Development

**Creating a new component:**

1. Create directory: `components/<name>/`
2. Create `types.go` for Go types and variant classes
3. Create `<name>.templ` for the component template
4. Use `dark:` prefixes for dark mode support

**Example from `components/button/types.go`:**
```go
func (v Variant) Classes() string {
    switch v {
    case Primary:
        return "bg-primary text-on-primary border-primary dark:bg-primary-dark dark:text-on-primary-dark dark:border-primary-dark"
    // ... other variants
    }
}
```

### Demo Pages

**Adding a component demo:**

1. Create file: `internal/pages/demo/components/<name>.templ`
2. Implement `DemoFragment` function returning `templ.Component`
3. Update sidebar in `layout.templ` with navigation link

**Demo page structure:**
- Header with component title and description
- Size controls (sm, md, lg, xl)
- Disabled toggle
- Split view or tab view with Original and GoATTH columns
- Code preview sections

## Testing

### E2E Tests

Located in `tests/e2e/`:

```bash
# Run all E2E tests
cd tests/e2e && go test -v

# Run specific test
cd tests/e2e && go test -v -run TestButtonVariants
```

**Test structure:**
- Uses Playwright for browser automation
- Tests visual parity between Original and GoATTH
- Validates color values and CSS classes
- Checks dark mode functionality

### Manual Testing Checklist

- [ ] Light mode renders correctly
- [ ] Dark mode toggles both Original and GoATTH simultaneously
- [ ] Theme selector changes color palette
- [ ] Size controls (sm, md, lg, xl) work
- [ ] Disabled toggle disables all buttons
- [ ] Code copy buttons work
- [ ] Responsive design works on mobile

## Design Guidelines

### Visual Parity Requirements

- **Colors:** Must match PenguinUI exactly (use browser dev tools to verify)
- **Spacing:** Use `gap-8` and `p-8` for grid layouts
- **Border Radius:** Respect theme's `--radius-radius` variable
- **Typography:** Use `font-title` for headings, default for body

### Component Classes Pattern

Always use both light and dark variants:
```css
bg-surface text-on-surface dark:bg-surface-dark dark:text-on-surface-dark
```

### Layout Grid

Standard grid for button previews:
```css
grid grid-cols-2 lg:grid-cols-4 gap-8 p-8 place-content-evenly place-items-center
```

## Common Issues

### Dark Mode Not Working

1. Check that `@custom-variant dark (&:where(.dark, .dark *))` exists in CSS
2. Verify `.dark` class is being added to `<html>` element
3. Ensure components use `dark:` prefixes

### Theme Not Switching

1. Check that `data-theme` attribute is set on `<html>`
2. Verify theme definitions exist in `all-themes.css`
3. Check CSS custom property values in browser dev tools

### Changes Not Appearing

1. Run `make generate` after editing `.templ` files
2. Clear browser cache
3. Restart server

## Server Ports

- **Default:** 7070
- **Configurable:** Use `-port` flag: `go run cmd/server/main.go -port 8080`

## GitHub Repository

https://github.com/guilycst/GoATTH-penguinui

## License

Same as PenguinUI (check original repository for license details)

---

**Last Updated:** March 2026
**Maintainer:** Agentic Coding Agent
