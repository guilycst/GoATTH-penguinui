package sidebar

import "github.com/a-h/templ"

// Item represents a single navigation item in the sidebar
type Item struct {
	// ID is the unique identifier for the item
	ID string
	// Label is the display text
	Label string
	// Icon is an optional icon (templ.Component)
	Icon templ.Component
	// Href is the link URL
	Href string
	// Active indicates if this item is currently active
	Active bool
	// Disabled prevents interaction
	Disabled bool
	// Badge is an optional badge text (e.g., "New", "Soon")
	Badge string
	// Items contains child items for nested navigation
	Items []Item
}

// Config holds configuration for the sidebar
type Config struct {
	// Items are the navigation items
	Items []Item
	// Logo is the logo component (optional)
	Logo templ.Component
	// LogoText is the text-only logo if no component provided
	LogoText string
	// LogoHref is the link for the logo (default: "/")
	LogoHref string
	// ShowSearch enables the search input
	ShowSearch bool
	// SearchPlaceholder is the placeholder text for search
	SearchPlaceholder string
	// Class allows additional CSS classes
	Class string
}

// Section represents a group of navigation items
type Section struct {
	// Title is the section header
	Title string
	// Items are the navigation items in this section
	Items []Item
}

// ContainerClasses returns the container CSS classes
func (cfg Config) ContainerClasses() string {
	return "fixed inset-y-0 left-0 z-40 w-64 transform border-r border-outline bg-surface transition-transform duration-200 ease-in-out lg:static lg:translate-x-0 dark:border-outline-dark dark:bg-surface-dark flex flex-col"
}

// NavClasses returns the navigation container classes
func (cfg Config) NavClasses() string {
	return "flex-1 overflow-y-auto sidebar-scroll p-4"
}
