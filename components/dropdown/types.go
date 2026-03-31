package dropdown

import "github.com/a-h/templ"

// TriggerMode determines how the dropdown is activated
type TriggerMode string

const (
	TriggerClick   TriggerMode = "click"
	TriggerHover   TriggerMode = "hover"
	TriggerContext TriggerMode = "context"
)

// Item represents a single menu item in the dropdown
type Item struct {
	// Label is the display text for the menu item
	Label string
	// Href is the link URL (use "#" for non-navigating items)
	Href string
	// Icon is an optional icon component rendered before the label
	Icon templ.Component
	// Shortcut is an optional keyboard shortcut label (e.g., "Z", "X")
	Shortcut string
	// ShortcutIcon is an optional icon for the shortcut modifier key
	ShortcutIcon templ.Component
}

// Section groups items with an optional heading.
type Section struct {
	// Items is the list of items in this section.
	Items []Item
}

// Config holds configuration for the dropdown component
type Config struct {
	// ID is the unique identifier for the dropdown
	ID string
	// Label is the text shown on the trigger button
	Label string
	// TriggerMode determines how the dropdown opens (click, hover, context)
	TriggerMode TriggerMode
	// Sections groups items with dividers between sections
	Sections []Section
	// TriggerIcon is an optional custom trigger icon (used for context menus)
	TriggerIcon templ.Component
}

// GetTriggerMode returns the trigger mode with a default of click
func (cfg Config) GetTriggerMode() TriggerMode {
	if cfg.TriggerMode == "" {
		return TriggerClick
	}
	return cfg.TriggerMode
}

// HasDividers returns true if there are multiple sections
func (cfg Config) HasDividers() bool {
	return len(cfg.Sections) > 1
}

// HasIcons returns true if any item has an icon
func (cfg Config) HasIcons() bool {
	for _, section := range cfg.Sections {
		for _, item := range section.Items {
			if item.Icon != nil {
				return true
			}
		}
	}
	return false
}

// HasShortcuts returns true if any item has a shortcut
func (cfg Config) HasShortcuts() bool {
	for _, section := range cfg.Sections {
		for _, item := range section.Items {
			if item.Shortcut != "" {
				return true
			}
		}
	}
	return false
}

// IsContextMenu returns true if this is a context menu trigger
func (cfg Config) IsContextMenu() bool {
	return cfg.GetTriggerMode() == TriggerContext
}

// MenuClasses returns the CSS classes for the dropdown menu container
func (cfg Config) MenuClasses() string {
	base := "absolute left-0 z-30 flex w-fit min-w-48 flex-col overflow-hidden rounded-radius border border-outline bg-surface-alt shadow-md dark:border-outline-dark dark:bg-surface-dark-alt"
	if cfg.IsContextMenu() {
		return "top-8 " + base
	}
	return "top-11 " + base
}

// ItemClasses returns the CSS classes for a dropdown menu item
func (cfg Config) ItemClasses(hasIcon bool) string {
	base := "bg-surface-alt px-4 py-2 text-sm text-on-surface hover:bg-surface-dark-alt/5 hover:text-on-surface-strong focus-visible:bg-surface-dark-alt/10 focus-visible:text-on-surface-strong focus-visible:outline-hidden dark:bg-surface-dark-alt dark:text-on-surface-dark dark:hover:bg-surface-alt/5 dark:hover:text-on-surface-dark-strong dark:focus-visible:bg-surface-alt/10 dark:focus-visible:text-on-surface-dark-strong"
	if hasIcon {
		return "flex items-center gap-2 " + base
	}
	return base
}

// ButtonClasses returns the CSS classes for the trigger button
func (cfg Config) ButtonClasses() string {
	if cfg.IsContextMenu() {
		return "inline-flex items-center bg-transparent transition hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-outline-strong active:opacity-100 dark:focus-visible:outline-outline-dark-strong"
	}
	return "inline-flex items-center gap-2 whitespace-nowrap rounded-radius border border-outline bg-surface-alt px-4 py-2 text-sm font-medium tracking-wide transition hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-outline-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:focus-visible:outline-outline-dark-strong"
}
