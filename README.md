# GoATTH PenguinUI

**GoATTH**: Go + Templ + Tailwind CSS + HTMX + Alpine.js

## About This Fork

This is a hard fork of [Penguin UI](https://www.penguinui.com) by Salar Houshvand, transformed from static HTML/Alpine.js components into a complete Go web component library.

### What's Changed?

| Original | GoATTH Fork |
|----------|-------------|
| Static HTML | Go + Templ templates |
| CDN assets | Configurable (CDN/Embedded/Custom) |
| Copy-paste | `go get` importable |
| Alpine.js only | HTMX + Alpine.js + Go backend |

## Credits

- **Original**: [Penguin UI](https://www.penguinui.com) by [Salar Houshvand](https://x.com/salar_houshvand)
- **License**: MIT (preserved from original)

## Project Structure

```
GoATTH-penguinui/
├── cmd/server/              # Demo server
├── components/              # GoATTH component library
│   └── button/             # Button component
│       ├── types.go        # Configuration types
│       └── button.templ    # Templ component
├── internal/
│   └── pages/demo/         # Demo pages
├── assets/css/             # Styles
├── tests/e2e/              # Playwright E2E tests
└── buttons/                # Original PenguinUI (preserved)
```

## Running the Demo

### Quick Start (Recommended: Air with Live Reload)

```bash
# Install dependencies (including Air)
make install
make install-air

# Run with live reload (auto-rebuilds on file changes)
make dev-air

# Server will start on http://localhost:8090
# Accordion Demo: http://localhost:8090/components/accordion
# Button Demo: http://localhost:8090/components/button
```

### Standard Development

```bash
# Install dependencies
make install

# Run the demo server
make dev
# or
go run cmd/server/main.go

# Server will start on http://localhost:8090
# - Original PenguinUI: http://localhost:8090/original/
# - GoATTH Components: http://localhost:8090/gottha/
```

## Running E2E Tests

E2E tests use [playwright-go](https://github.com/playwright-community/playwright-go) (following the tks-console pattern) for Go-based browser automation.

```bash
# Using just (from repo root)
just gp-test-e2e                    # Run all E2E tests
just gp-test-e2e-one TestButton     # Run specific test

# Or directly
go test ./GoATTH-penguinui/tests/e2e/... -v

# First time setup - install Playwright browsers
just gp-install-playwright
# or
go install github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1
playwright install chromium
```

### Test Results

Tests automatically:
- Start the demo server
- Run browser automation tests
- Capture screenshots on failures to `test-results/screenshots/`
- Verify both Original PenguinUI and GoATTH component rendering

### Current Test Coverage

- **Button Component**: Verifies all 8 variants render correctly, HTMX attributes, Alpine.js integration
- **Screenshots**: Auto-captured for visual debugging

## Component Usage

### Button Component

```go
import "github.com/guilycst/GoATTH-penguinui/components/button"

// Basic button
@button.Button(button.Config{
    Variant: button.Primary,
    Type:    "button",
}) {
    Click Me
}

// With HTMX
@button.Button(button.Config{
    Variant: button.Primary,
    HTMX: &button.HTMXConfig{
        Post:   "/api/action",
        Target: "#result",
        Swap:   "innerHTML",
    },
}) {
    Submit
}

// With Alpine.js
@button.Button(button.Config{
    Variant: button.Primary,
    Alpine: &button.AlpineConfig{
        OnClick: "modalIsOpen = true",
    },
}) {
    Open Modal
}
```

## Development

### Building

```bash
make build
```

### Generating Templ Files

```bash
make generate
```

### Testing

```bash
# Go tests
make test

# E2E tests
make test-e2e
```

## License

MIT License - See original [Penguin UI](https://www.penguinui.com/docs/license) for details.
