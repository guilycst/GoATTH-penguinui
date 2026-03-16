# GoATTH Components - Usage Guide

This guide explains how to use GoATTH (Go + Alpine.js + Tailwind CSS + HTMX + Templ) components from the tks-console and other projects.

## Installation

### 1. Add Dependency

In your `go.mod`:

```go
require github.com/guilycst/GoATTH-penguinui v0.1.0
```

Or install via go get:

```bash
go get github.com/guilycst/GoATTH-penguinui@latest
```

### 2. Install Tailwind CSS Dependencies

Your project must have Tailwind CSS v4+ configured with the same theme variables as GoATTH:

**Required CSS files:**
- Copy `all-themes.css` from GoATTH to your project's CSS directory
- Import it in your main CSS file

```css
/* your-project/css/main.css */
@import "tailwindcss";
@import "../all-themes.css";
```

### 3. Required JavaScript

Include these CDN links in your HTML head:

```html
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
<script src="https://unpkg.com/htmx.org@2.0.8/dist/htmx.min.js"></script>
```

For components using Alpine.js collapse plugin (like Accordion):
```html
<script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js"></script>
```

## Available Components

### Accordion

A collapsible accordion component with multiple variants.

**Import:**
```go
import "github.com/guilycst/GoATTH-penguinui/components/accordion"
```

**Basic Usage:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    Items: []accordion.AccordionItem{
        {
            ID:      "section1",
            Title:   "Section 1",
            Content: templ.Raw("<p>This is the content for section 1.</p>"),
        },
        {
            ID:      "section2", 
            Title:   "Section 2",
            Content: templ.Raw("<p>This is the content for section 2.</p>"),
        },
    },
})
```

**Allow Multiple Open:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    AllowMultiple: true,
    Items: []accordion.AccordionItem{
        {ID: "item1", Title: "Item 1", Content: content1},
        {ID: "item2", Title: "Item 2", Content: content2},
    },
})
```

**No Background Variant:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    Variant: accordion.NoBackground,
    Items:   items,
})
```

**With Icons:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    Items: []accordion.AccordionItem{
        {
            ID:      "settings",
            Title:   "Settings",
            Icon:    icons.SettingsIcon(), // Your icon component
            Content: SettingsContent(),
        },
    },
})
```

**Initially Expanded:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    Items: []accordion.AccordionItem{
        {
            ID:                "details",
            Title:             "Details",
            Content:           DetailsContent(),
            InitiallyExpanded: true,
        },
    },
})
```

**Disabled Item:**
```go
@accordion.Accordion(accordion.AccordionConfig{
    Items: []accordion.AccordionItem{
        {
            ID:       "restricted",
            Title:    "Restricted Section",
            Content:  Content(),
            Disabled: true,
        },
    },
})
```

### Button

A versatile button component with multiple variants.

**Import:**
```go
import "github.com/guilycst/GoATTH-penguinui/components/button"
```

**Basic Usage:**
```go
@button.Button(button.Config{
    Variant: button.Primary,
    Type:    "button",
}) {
    Click Me
}
```

**With HTMX:**
```go
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
```

**With Alpine.js:**
```go
@button.Button(button.Config{
    Variant: button.Primary,
    Alpine: &button.AlpineConfig{
        OnClick: "modalIsOpen = true",
    },
}) {
    Open Modal
}
```

**All Variants:**
- `button.Primary` - Black background, white text
- `button.Secondary` - Dark gray, white text
- `button.Alternate` - Light gray, dark text
- `button.Inverse` - Inverted colors
- `button.Info` - Sky blue, white text
- `button.Danger` - Red, white text
- `button.Warning` - Yellow/amber, dark text
- `button.Success` - Green, dark text

**Sizes:**
```go
@button.Button(button.Config{
    Variant: button.Primary,
    Size:    button.SizeSmall,  // xs text
})

@button.Button(button.Config{
    Variant: button.Primary,
    Size:    button.SizeMedium, // sm text (default)
})

@button.Button(button.Config{
    Variant: button.Primary,
    Size:    button.SizeLarge,  // base text
})

@button.Button(button.Config{
    Variant: button.Primary,
    Size:    button.SizeXLarge, // lg text
})
```

## Configuration Types

### AccordionItem

```go
type AccordionItem struct {
    ID                string          // Unique identifier
    Title             string          // Header text
    Content           templ.Component // Body content
    Icon              templ.Component // Optional icon
    Disabled          bool            // Disable interaction
    InitiallyExpanded bool            // Start expanded
}
```

### AccordionConfig

```go
type AccordionConfig struct {
    Items         []AccordionItem  // Accordion sections
    AllowMultiple bool             // Multiple open at once
    Variant       Variant          // Visual style
    ID            string           // Container ID
    Class         string           // Additional CSS classes
}
```

### Button Config

```go
type Config struct {
    Variant     Variant       // Button style
    Size        Size          // Button size
    Type        string        // HTML type attribute
    Disabled    bool          // Disabled state
    ID          string        // Element ID
    Class       string        // Additional classes
    HTMX        *HTMXConfig   // HTMX attributes
    Alpine      *AlpineConfig // Alpine.js directives
    LoadingText string        // Loading state text
}
```

## Theming

### Available Themes

GoATTH supports all 13 PenguinUI themes:

1. **Minimal** (default) - Black/white, no border radius
2. **Arctic** - Cool blue tones
3. **Modern** - Clean professional look
4. **High-Contrast** - Maximum accessibility
5. **And more...**

### Switching Themes

Set the theme via data attribute on `<html>`:

```html
<html data-theme="modern">
```

Or with JavaScript:

```javascript
document.documentElement.setAttribute('data-theme', 'modern');
```

### Dark Mode

Add/remove the `dark` class on `<html>`:

```javascript
document.documentElement.classList.toggle('dark');
```

## Best Practices

### 1. Content Components

Create separate templ components for accordion content to keep code clean:

```go
templ SettingsContent() {
    <div class="space-y-4">
        <div>
            <label class="block text-sm font-medium">Name</label>
            <input type="text" class="mt-1 block w-full" />
        </div>
        <div>
            <label class="block text-sm font-medium">Email</label>
            <input type="email" class="mt-1 block w-full" />
        </div>
    </div>
}

// Use it
@accordion.Accordion(accordion.AccordionConfig{
    Items: []accordion.AccordionItem{
        {Title: "Settings", Content: SettingsContent()},
    },
})
```

### 2. Icons

Pass icons as templ.Components:

```go
func InfoIcon() templ.Component {
    return templ.Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="size-5" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
    </svg>`)
}
```

### 3. HTMX Integration

Components work seamlessly with HTMX for dynamic updates:

```go
// Initial render
@accordion.Accordion(accordion.AccordionConfig{
    ID: "cluster-accordion",
    Items: []accordion.AccordionItem{
        {
            ID:      "node-pools",
            Title:   "Node Pools",
            Content: NodePoolsTable(cluster.NodePools),
        },
    },
})

// Update fragment via HTMX
func HandleNodePoolsUpdate(w http.ResponseWriter, r *http.Request) {
    clusterID := r.URL.Query().Get("cluster_id")
    nodePools := fetchNodePools(clusterID)
    
    // Render just the content that changed
    accordion.AccordionItem{
        ID:      "node-pools",
        Title:   "Node Pools",
        Content: NodePoolsTable(nodePools),
    }.Render(r.Context(), w)
}
```

### 4. Testing

Always test your component implementations:

```go
// In your project tests
func TestAccordion_Integration(t *testing.T) {
    // Use GoATTH's visual testing utilities
    // See GoATTH tests for examples
}
```

## Troubleshooting

### Component not styled correctly?

1. Check that Tailwind CSS is processing the component files
2. Verify `all-themes.css` is imported
3. Ensure the `data-theme` attribute is set on `<html>`

### Alpine.js not working?

1. Verify Alpine.js is loaded before components render
2. Check browser console for Alpine.js errors
3. For collapse animations, ensure `@alpinejs/collapse` is loaded

### Dark mode not working?

1. Add `dark` class to `<html>` element
2. Verify CSS custom properties are defined in `all-themes.css`
3. Check that `dark:` prefixes are in component classes

## Examples

See the `/components` directory in the GoATTH repository for complete examples of each component with visual parity tests.

Run the demo server:

```bash
cd /path/to/GoATTH-penguinui
go run cmd/server/main.go -port 8090
```

Then visit:
- http://localhost:8090/components/button
- http://localhost:8090/components/accordion

## Contributing

To add new components:

1. Create component in `components/<name>/`
2. Copy original HTML to `fixtures/`
3. Create demo page in `internal/pages/demo/components/`
4. Write E2E tests with visual parity checks
5. Document in this guide

## License

MIT License (same as PenguinUI)
