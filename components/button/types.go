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
	SizeXLarge Size = "xl"
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
	classes := "whitespace-nowrap rounded-2xl font-medium tracking-wide transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed"

	// Size classes
	switch cfg.Size {
	case SizeSmall:
		classes += " px-4 py-2 text-xs"
	case SizeLarge:
		classes += " px-4 py-2 text-base"
	case SizeXLarge:
		classes += " px-4 py-2 text-lg"
	default: // Medium
		classes += " px-4 py-2 text-sm"
	}

	// Variant-specific classes with inline styles for colors
	return classes
}

// variantStyles returns the inline styles for a variant
func variantStyles(variant Variant) string {
	switch variant {
	case Primary:
		return "background-color: #000000; color: #f8fafc; border: 1px solid #000000;"
	case Secondary:
		return "background-color: #1e293b; color: #ffffff; border: 1px solid #1e293b;"
	case Alternate:
		return "background-color: #f1f5f9; color: #0f172a; border: 1px solid #f1f5f9;"
	case Inverse:
		return "background-color: #000000; color: #f8fafc; border: 1px solid #000000;"
	case Info:
		return "background-color: #0ea5e9; color: #ffffff; border: 1px solid #0ea5e9;"
	case Danger:
		return "background-color: #ef4444; color: #ffffff; border: 1px solid #ef4444;"
	case Warning:
		return "background-color: #fcd34d; color: #78350f; border: 1px solid #fcd34d;"
	case Success:
		return "background-color: #86efac; color: #0f172a; border: 1px solid #86efac;"
	default:
		return ""
	}
}

// variantOutlineStyles returns the outline color for focus states
func variantOutlineStyles(variant Variant) string {
	switch variant {
	case Primary:
		return "#000000"
	case Secondary:
		return "#1e293b"
	case Alternate:
		return "#f1f5f9"
	case Inverse:
		return "#000000"
	case Info:
		return "#0ea5e9"
	case Danger:
		return "#ef4444"
	case Warning:
		return "#fcd34d"
	case Success:
		return "#86efac"
	default:
		return "#000000"
	}
}
