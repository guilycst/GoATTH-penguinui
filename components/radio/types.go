package radio

// Variant represents radio button color variants
type Variant string

const (
	Primary   Variant = "primary"
	Secondary Variant = "secondary"
	Info      Variant = "info"
	Success   Variant = "success"
	Warning   Variant = "warning"
	Danger    Variant = "danger"
)

// Config holds configuration for a single radio button
type Config struct {
	// ID is the unique identifier for the radio input
	ID string
	// Name is the form field name (shared by all radios in a group)
	Name string
	// Value is the form field value
	Value string
	// Label is the text displayed next to the radio
	Label string
	// Checked sets the initial checked state
	Checked bool
	// Disabled disables the radio
	Disabled bool
	// Variant determines the color scheme (default: Primary)
	Variant Variant
	// Description adds helper text below the label
	Description string
	// BadgeColor wraps the label in a semi-solid badge. Accepts: "success", "danger", "warning", "info", "neutral", "primary", "secondary"
	BadgeColor string
	// Container wraps the radio in a bordered container
	Container bool
}

// GroupConfig holds configuration for a group of radio buttons
type GroupConfig struct {
	// Title is the heading above the group
	Title string
	// Items are the radio buttons in the group
	Items []Config
}

// checkedBorderClass returns the checked border color class
func (cfg Config) checkedBorderClass() string {
	switch cfg.Variant {
	case Secondary:
		return "checked:border-secondary dark:checked:border-secondary-dark"
	case Info:
		return "checked:border-info dark:checked:border-info"
	case Success:
		return "checked:border-success dark:checked:border-success"
	case Warning:
		return "checked:border-warning dark:checked:border-warning"
	case Danger:
		return "checked:border-danger dark:checked:border-danger"
	default:
		return "checked:border-primary dark:checked:border-primary-dark"
	}
}

// checkedBgClass returns the checked background color class
func (cfg Config) checkedBgClass() string {
	switch cfg.Variant {
	case Secondary:
		return "checked:before:bg-secondary dark:checked:before:bg-secondary-dark"
	case Info:
		return "checked:before:bg-info dark:checked:before:bg-info"
	case Success:
		return "checked:before:bg-success dark:checked:before:bg-success"
	case Warning:
		return "checked:before:bg-warning dark:checked:before:bg-warning"
	case Danger:
		return "checked:before:bg-danger dark:checked:before:bg-danger"
	default:
		return "checked:before:bg-primary dark:checked:before:bg-primary-dark"
	}
}

// focusCheckedClass returns the focus outline color when checked
func (cfg Config) focusCheckedClass() string {
	switch cfg.Variant {
	case Secondary:
		return "checked:focus:outline-secondary dark:checked:focus:outline-secondary-dark"
	case Info:
		return "checked:focus:outline-info dark:checked:focus:outline-info"
	case Success:
		return "checked:focus:outline-success dark:checked:focus:outline-success"
	case Warning:
		return "checked:focus:outline-warning dark:checked:focus:outline-warning"
	case Danger:
		return "checked:focus:outline-danger dark:checked:focus:outline-danger"
	default:
		return "checked:focus:outline-primary dark:checked:focus:outline-primary-dark"
	}
}

// BadgeClasses returns CSS classes for a badge-styled label
func BadgeClasses(color string) string {
	base := "w-fit rounded-radius px-2 py-0.5 text-xs font-medium"
	switch color {
	case "success":
		return base + " bg-success/10 text-success"
	case "danger":
		return base + " bg-danger/10 text-danger"
	case "warning":
		return base + " bg-warning/10 text-warning"
	case "info":
		return base + " bg-info/10 text-info"
	case "primary":
		return base + " bg-primary/10 text-primary dark:text-primary-dark"
	case "secondary":
		return base + " bg-secondary/10 text-secondary"
	case "neutral":
		return base + " bg-on-surface/10 text-on-surface dark:text-on-surface-dark"
	default:
		return ""
	}
}

// InputClasses returns the full CSS class string for the radio input
func (cfg Config) InputClasses() string {
	base := "before:content[''] peer relative size-4 shrink-0 appearance-none overflow-hidden rounded-full border border-outline bg-surface-alt before:absolute before:inset-0 before:scale-0 before:rounded-full before:transition before:duration-200 checked:before:scale-[0.55] focus:outline-2 focus:outline-offset-2 focus:outline-outline-strong active:outline-offset-0 disabled:cursor-not-allowed dark:border-outline-dark dark:bg-surface-dark-alt dark:focus:outline-outline-dark-strong"

	if cfg.Container {
		base = "before:content[''] peer relative size-4 shrink-0 appearance-none overflow-hidden rounded-full border border-outline bg-surface before:absolute before:inset-0 before:scale-0 before:rounded-full before:transition before:duration-200 checked:before:scale-[0.55] focus:outline-2 focus:outline-offset-2 focus:outline-outline-strong active:outline-offset-0 disabled:cursor-not-allowed dark:border-outline-dark dark:bg-surface-dark dark:focus:outline-outline-dark-strong"
	}

	return base + " " + cfg.checkedBorderClass() + " " + cfg.checkedBgClass() + " " + cfg.focusCheckedClass()
}
