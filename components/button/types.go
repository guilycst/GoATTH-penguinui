package button

// Variant represents button style variants
type Variant string

const (
	Primary   Variant = "primary"
	Secondary Variant = "secondary"
	Alternate Variant = "alternate"
	Inverse   Variant = "inverse"
	Info      Variant = "info"
	Danger    Variant = "danger"
	Warning   Variant = "warning"
	Success   Variant = "success"
)

// Size represents button sizes
type Size string

const (
	SizeSmall  Size = "sm"
	SizeMedium Size = "md"
	SizeLarge  Size = "lg"
)

// HTMXConfig holds HTMX attributes for server-side interactions
type HTMXConfig struct {
	Get       string
	Post      string
	Put       string
	Delete    string
	Patch     string
	Target    string
	Swap      string
	Trigger   string
	Indicator string
	PushURL   bool
	Confirm   string
	Vals      string
}

// AlpineConfig holds Alpine.js directives for client-side interactions
type AlpineConfig struct {
	OnClick      string
	BindDisabled string
	Show         string
	Transition   bool
	Data         string
}

// Config holds all configuration for a Button component
type Config struct {
	Variant     Variant
	Size        Size
	Type        string
	Disabled    bool
	ID          string
	Class       string
	HTMX        *HTMXConfig
	Alpine      *AlpineConfig
	LoadingText string
}

// buttonClasses returns the CSS classes for a button based on its configuration
func buttonClasses(cfg Config) string {
	// Base classes
	classes := "whitespace-nowrap rounded-radius font-medium tracking-wide transition text-center focus-visible:outline-2 focus-visible:outline-offset-2 active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed"

	// Variant-specific classes (matching PenguinUI exactly)
	switch cfg.Variant {
	case Primary:
		classes += " bg-primary border border-primary text-on-primary hover:opacity-75 focus-visible:outline-primary dark:bg-primary-dark dark:border-primary-dark dark:text-on-primary-dark dark:focus-visible:outline-primary-dark"
	case Secondary:
		classes += " bg-secondary border border-secondary text-on-secondary hover:opacity-75 focus-visible:outline-secondary dark:bg-secondary-dark dark:border-secondary-dark dark:text-on-secondary-dark dark:focus-visible:outline-secondary-dark"
	case Alternate:
		classes += " bg-surface-alt border border-surface-alt text-on-surface-strong hover:opacity-75 focus-visible:outline-surface-alt dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark-strong dark:focus-visible:outline-surface-dark-alt"
	case Inverse:
		classes += " bg-surface-dark border border-surface-dark text-on-surface-dark hover:opacity-75 focus-visible:outline-surface-dark dark:bg-surface dark:border-surface dark:text-on-surface dark:focus-visible:outline-surface"
	case Info:
		classes += " bg-info border border-info text-onInfo hover:opacity-75 focus-visible:outline-info dark:bg-info dark:border-info dark:text-onInfo dark:focus-visible:outline-info"
	case Danger:
		classes += " bg-danger border border-danger text-onDanger hover:opacity-75 focus-visible:outline-danger dark:bg-danger dark:border-danger dark:text-onDanger dark:focus-visible:outline-danger"
	case Warning:
		classes += " bg-warning border border-warning text-onWarning hover:opacity-75 focus-visible:outline-warning dark:bg-warning dark:border-warning dark:text-onWarning dark:focus-visible:outline-warning"
	case Success:
		classes += " bg-success border border-success text-onSuccess hover:opacity-75 focus-visible:outline-success dark:bg-success dark:border-success dark:text-onSuccess dark:focus-visible:outline-success"
	}

	// Size modifiers
	switch cfg.Size {
	case SizeSmall:
		classes += " px-3 py-1.5 text-xs"
	case SizeLarge:
		classes += " px-6 py-3 text-base"
	default: // Medium
		classes += " px-4 py-2 text-sm"
	}

	return classes
}
