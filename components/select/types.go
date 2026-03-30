package selectfield

import "github.com/a-h/templ"

// State represents the validation state of the select
type State string

const (
	StateDefault State = ""
	StateError   State = "error"
	StateSuccess State = "success"
)

// Option represents a single selectable item
type Option struct {
	Value    string
	Label    string
	Selected bool
}

// Config holds configuration for the select component
type Config struct {
	// ID is a unique identifier for the select element
	ID string
	// Name is the form field name
	Name string
	// Label is the label text shown above the select
	Label string
	// Placeholder text when no selection (default: "Please Select")
	Placeholder string
	// Options is the list of available options
	Options []Option
	// State is the validation state (error, success, or default)
	State State
	// HelperText is shown below the select (e.g., error or success message)
	HelperText string
	// Disabled disables the select
	Disabled bool
	// Autocomplete sets the autocomplete attribute
	Autocomplete string
	// Class allows additional CSS classes on the wrapper
	Class string
	// AlpineModel sets x-model on the select for Alpine.js binding
	AlpineModel string
	// AlpineBindDisabled sets x-bind:disabled on the select
	AlpineBindDisabled string
	// Readonly renders the select as disabled (grayed out) + hidden input with value so it still submits
	Readonly bool
	// Attrs allows arbitrary HTML attributes on the <select> element (e.g., hx-post, hx-indicator)
	Attrs templ.Attributes
}

// ContainerClasses returns CSS classes for the outer wrapper.
// When Label is set (standalone usage), max-w-xs is applied.
// When Label is empty (inside form FieldGroup), full width.
func (cfg Config) ContainerClasses() string {
	base := "relative flex w-full flex-col gap-1 text-on-surface dark:text-on-surface-dark"
	if cfg.Label != "" {
		base += " max-w-xs"
	}
	if cfg.Class != "" {
		base += " " + cfg.Class
	}
	return base
}

// SelectClasses returns CSS classes for the select element
func (cfg Config) SelectClasses() string {
	base := "w-full appearance-none rounded-radius border bg-surface-alt px-4 py-2 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-not-allowed disabled:opacity-75 dark:bg-surface-dark-alt/50 dark:focus-visible:outline-primary-dark"

	switch cfg.State {
	case StateError:
		return base + " border-danger"
	case StateSuccess:
		return base + " border-success"
	default:
		return base + " border-outline dark:border-outline-dark"
	}
}

// LabelClasses returns CSS classes for the label
func (cfg Config) LabelClasses() string {
	base := "w-fit pl-0.5 text-sm"

	switch cfg.State {
	case StateError:
		return "flex w-fit gap-1 pl-0.5 text-sm text-danger"
	case StateSuccess:
		return "flex w-fit gap-1 pl-0.5 text-sm text-success"
	default:
		return base
	}
}

// GetPlaceholder returns the placeholder text
func (cfg Config) GetPlaceholder() string {
	if cfg.Placeholder != "" {
		return cfg.Placeholder
	}
	return "Please Select"
}

// SelectedValue returns the value of the first selected option, or empty string
func (cfg Config) SelectedValue() string {
	for _, opt := range cfg.Options {
		if opt.Selected {
			return opt.Value
		}
	}
	return ""
}

// IsEffectivelyDisabled returns true if the select should render as disabled (Disabled or Readonly)
func (cfg Config) IsEffectivelyDisabled() bool {
	return cfg.Disabled || cfg.Readonly
}
