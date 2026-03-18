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

## Current Status

### Completed Components ✅
1. **Button** - 8 variants, HTMX support, Alpine.js integration
2. **Accordion** - Static & server-loaded content, multiple variants, visual tests
3. **Sidebar** - Navigation component with collapsible sections
4. **Avatar** - Image, initials, icon placeholders, status indicators, 5 sizes
5. **Badge** - Solid & soft variants, icons, indicators, notification badges, 3 sizes
6. **Banner** - Dismissible, CTA buttons, cookie consent, color variants
7. **Card** - Vertical/horizontal layouts, product card, pricing card, testimonial
8. **Combobox** - Single/multi-select, search, images, pre-selected values, disabled state

### Pending Components ⏳
- Modal
- Alert
- Input
- Select
- Checkbox
- Radio
- Toggle
- Toast/Notification
- Table
- Tabs
- Pagination
- Progress
- Skeleton
- Spinner
- Tooltip
- Dropdown
- Date Picker
- File Input
- Text Area
- Range Slider
- Rating/Stars
- Steps/Wizard
- Timeline
- Chat/Message
- Calendar
- Charts (optional)

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

### Component Development Workflow

The established workflow for creating new components:

1. **Reference Analysis**
   - Read existing PenguinUI HTML reference files in `/combobox/`, `/accordion/`, etc.
   - Check tks-console implementations for feature ideas (optional)
   - Document visual design requirements (CSS classes, colors, spacing)

2. **Create Component Files**
   ```
   components/<name>/
   ├── <name>.templ    # Component template
   ├── types.go        # Configuration types
   └── <name>_templ.go # Generated (via make generate)
   ```

3. **Create Demo Page**
   ```
   internal/pages/demo/components/<name>.templ
   ```

4. **Update Server Routes**
   - Add route handler in `internal/server/server.go`
   - Add sidebar item in `internal/pages/demo/layout.templ`

5. **Create E2E Tests**
   ```
   tests/e2e/<name>_test.go
   ```

6. **Generate and Test**
   ```bash
   make generate && make dev-air
   # Test at http://localhost:8090/components/<name>
   ```

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

## Lessons Learned & Best Practices

### Alpine.js Integration

**Initialize arrays properly:**
Always initialize Alpine.js arrays with `[]` not `null`:
```go
// BAD: selectedValues: null
// GOOD: selectedValues: []

selectedValuesData := selectedJSON
if selectedValuesData == "null" {
    selectedValuesData = "[]"
}
```

**Computed properties need defensive checks:**
```javascript
get selectedOption() {
    if (!this.selectedValues || this.selectedValues.length === 0) return null;
    // ... rest of logic
}
```

**Template rendering timing:**
- Alpine.js renders options dynamically using `x-for`
- Options may not exist in DOM immediately on page load
- Wait for Alpine.js to initialize: `page.WaitForTimeout(800)` in E2E tests
- Use `WaitForSelector` with timeout instead of immediate assertions

### Component Structure Patterns

**Types file structure:**
```go
package combobox

// Config holds all component configuration
type Config struct {
    ID       string
    Label    string
    Options  []Option
    Selected []string
    // ... other fields
}

// Helper methods for CSS classes
func (cfg Config) TriggerClasses() string { ... }
func (cfg Config) DropdownClasses() string { ... }
```

**Template structure:**
```templ
// Single main component entry point
templ Combobox(cfg Config) {
    if cfg.IsMultiple() {
        @multiSelectCombobox(cfg)
    } else {
        @singleSelectCombobox(cfg)
    }
}

// Private helper templates
templ singleSelectCombobox(cfg Config) { ... }
templ multiSelectCombobox(cfg Config) { ... }
```

**Alpine.js data generation:**
```go
func singleSelectData(...) string {
    // Convert Go data to JSON for Alpine.js
    optsJSON, _ := json.Marshal(options)
    selectedJSON, _ := json.Marshal(selected)
    
    // Ensure arrays are never null
    if selectedJSON == "null" {
        selectedJSON = "[]"
    }
    
    return fmt.Sprintf(`{
        allOptions: %s,
        selectedValues: %s,
        // ...
    }`, optsJSON, selectedJSON)
}
```

### E2E Testing Best Practices

**Handle Alpine.js timing:**
```go
// Wait for Alpine.js to render
trigger.Click()
page.WaitForTimeout(800)  // Give Alpine.js time to update

// Check state after waiting
expanded, _ := trigger.GetAttribute("aria-expanded")
```

**Use page reloads for clean state:**
```go
// Each test gets a fresh page
_, err := page.Reload(playwright.PageReloadOptions{
    WaitUntil: playwright.WaitUntilStateNetworkidle,
})
```

**Test for DOM presence, not just visibility:**
```go
// Element may exist but be hidden by x-show
count, _ := dropdown.Count()
if count > 0 {
    t.Log("✓ Dropdown exists in DOM")
}
```

### Component Patterns from Reference

**Merging implementations:**
When multiple reference implementations exist (e.g., tks-console + PenguinUI):
1. Start with PenguinUI visual design (CSS classes, colors, spacing)
2. Add tks-console features (search, multi-select, lazy loading)
3. Ensure both work together visually

**Event dispatching:**
```javascript
// Dispatch custom events for parent components
this.$dispatch('combobox-change', { 
    id: this.id, 
    value: option.value,
    values: this.selectedValues
});
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

### Alpine.js Options Not Rendering

**Symptom:** Dropdown opens but options are empty

**Causes:**
1. `selectedValues` initialized as `null` instead of `[]`
2. `filteredOptions` getter fails due to null check
3. `x-for` template not rendering because data is invalid

**Fix:**
```go
// In Go code, ensure arrays are never null
selectedJSON, _ := json.Marshal(cfg.Selected)
if string(selectedJSON) == "null" {
    selectedJSON = []byte("[]")
}
```

## Server Ports

- **Default:** 7070
- **Configurable:** Use `-port` flag: `go run cmd/server/main.go -port 8080`

## GitHub Repository

https://github.com/guilycst/GoATTH-penguinui

## License

Same as PenguinUI (check original repository for license details)

---

**Last Updated:** March 17, 2026
**Status:** 8 components completed, 16+ pending
**Maintainer:** Agentic Coding Agent
