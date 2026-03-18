package table

import "github.com/a-h/templ"

// Variant represents table style variants
type Variant string

const (
	Default  Variant = "default"
	Striped  Variant = "striped"
	WithCheckbox Variant = "checkbox"
)

// Column defines a table column header
type Column struct {
	// Key is the column identifier used to look up cell values in Row.Cells
	Key string
	// Label is the display text for the column header
	Label string
}

// Cell holds the content for a single table cell
type Cell struct {
	// Text is plain text content
	Text string
	// Component is a templ component to render (overrides Text)
	Component templ.Component
}

// Row represents a single table row
type Row struct {
	// ID is a unique identifier for the row (used for checkbox IDs)
	ID string
	// Cells maps column keys to cell content
	Cells map[string]Cell
}

// Config holds configuration for the table component
type Config struct {
	// ID is the table element ID
	ID string
	// Columns defines the table headers
	Columns []Column
	// Rows holds the table data
	Rows []Row
	// Variant determines the table style
	Variant Variant
	// ShowCheckbox adds a select-all checkbox column
	ShowCheckbox bool
	// Class allows additional CSS classes on the container
	Class string
}

// ContainerClasses returns the outer wrapper CSS classes
func (cfg Config) ContainerClasses() string {
	base := "overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"
	if cfg.Class != "" {
		base += " " + cfg.Class
	}
	return base
}

// TableClasses returns the <table> element CSS classes
func (cfg Config) TableClasses() string {
	return "w-full text-left text-sm text-on-surface dark:text-on-surface-dark"
}

// TheadClasses returns the <thead> CSS classes
func (cfg Config) TheadClasses() string {
	return "border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"
}

// TbodyClasses returns the <tbody> CSS classes
func (cfg Config) TbodyClasses() string {
	return "divide-y divide-outline dark:divide-outline-dark"
}

// RowClasses returns CSS classes for a table row
func (cfg Config) RowClasses() string {
	if cfg.Variant == Striped {
		return "even:bg-primary/5 dark:even:bg-primary-dark/10"
	}
	return ""
}

// CellClasses returns CSS classes for a table cell
func (cfg Config) CellClasses() string {
	return "p-4"
}

// HeaderCellClasses returns CSS classes for a header cell
func (cfg Config) HeaderCellClasses() string {
	return "p-4"
}

// CheckboxClasses returns CSS classes for checkboxes
func (cfg Config) CheckboxClasses() string {
	return "before:content[''] peer relative size-4 appearance-none overflow-hidden rounded border border-outline bg-surface before:absolute before:inset-0 checked:border-primary checked:before:bg-primary focus:outline-2 focus:outline-offset-2 focus:outline-outline-strong checked:focus:outline-primary active:outline-offset-0 dark:border-outline-dark dark:bg-surface-dark-alt dark:checked:border-primary-dark dark:checked:before:bg-primary-dark dark:focus:outline-outline-dark-strong dark:checked:focus:outline-primary-dark"
}

// ActionButtonClasses returns CSS classes for action buttons in table cells
func ActionButtonClasses() string {
	return "cursor-pointer whitespace-nowrap rounded-radius bg-transparent p-0.5 font-semibold text-primary outline-primary hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 active:opacity-100 active:outline-offset-0 dark:text-primary-dark dark:outline-primary-dark"
}

// StatusBadgeClasses returns CSS classes for status badges
func StatusBadgeClasses(status string) string {
	base := "inline-flex overflow-hidden rounded-radius px-1 py-0.5 text-xs font-medium"
	switch status {
	case "active", "success":
		return base + " border-success text-success bg-success/10"
	case "canceled", "danger":
		return base + " border-danger text-danger bg-danger/10"
	default:
		return base
	}
}
