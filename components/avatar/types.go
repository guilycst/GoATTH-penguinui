package avatar

import "github.com/a-h/templ"

// Size represents avatar size variants
type Size string

const (
	SizeXS Size = "xs" // Extra small
	SizeSM Size = "sm" // Small
	SizeMD Size = "md" // Medium (default)
	SizeLG Size = "lg" // Large
	SizeXL Size = "xl" // Extra large
)

// Variant represents avatar style variants
type Variant string

const (
	Default   Variant = "default"
	Inverse   Variant = "inverse"
	Primary   Variant = "primary"
	Secondary Variant = "secondary"
	Info      Variant = "info"
	Success   Variant = "success"
	Warning   Variant = "warning"
	Danger    Variant = "danger"
)

// Shape represents avatar shape
type Shape string

const (
	ShapeCircle Shape = "circle" // Rounded full (default)
	ShapeSquare Shape = "square" // Rounded corners
)

// Status represents online/offline status indicator
type Status string

const (
	StatusOffline Status = "offline"
	StatusInfo    Status = "info"
	StatusSuccess Status = "success"
	StatusWarning Status = "warning"
	StatusDanger  Status = "danger"
)

// Config holds configuration for the avatar
type Config struct {
	// Src is the image URL (if using image avatar)
	Src string
	// Alt is the alt text for accessibility
	Alt string
	// Initials are displayed when no image (e.g., "JS")
	Initials string
	// Size of the avatar
	Size Size
	// Variant determines the color scheme (for initials/icon placeholders)
	Variant Variant
	// Shape of the avatar (circle or square)
	Shape Shape
	// Border adds a colored border (for image avatars)
	Border bool
	// BorderColor is the border color (defaults to variant color if empty)
	BorderColor string
	// Status adds a status indicator dot
	Status Status
	// Icon is an optional icon component (replaces initials)
	Icon templ.Component
	// Class allows additional CSS classes
	Class string
}

// SizeClasses returns the CSS classes for the size
func (cfg Config) SizeClasses() string {
	switch cfg.Size {
	case SizeXS:
		return "size-8 text-xs"
	case SizeSM:
		return "size-10 text-sm"
	case SizeMD:
		return "size-14 text-2xl"
	case SizeLG:
		return "size-20 text-3xl"
	case SizeXL:
		return "size-24 text-4xl"
	default:
		return "size-14 text-2xl"
	}
}

// ShapeClasses returns the CSS classes for the shape
func (cfg Config) ShapeClasses() string {
	switch cfg.Shape {
	case ShapeSquare:
		return "rounded-md"
	default:
		return "rounded-full"
	}
}

// VariantClasses returns the CSS classes for the variant (for initials/icon)
func (cfg Config) VariantClasses() string {
	switch cfg.Variant {
	case Inverse:
		return "border border-outline-dark bg-surface-dark-alt text-on-surface-dark/80 dark:border-outline dark:bg-surface-alt dark:text-on-surface/80"
	case Primary:
		return "border border-primary bg-primary text-on-primary/80 dark:border-primary-dark dark:bg-primary-dark dark:text-on-primary-dark/80"
	case Secondary:
		return "border border-secondary bg-secondary text-on-secondary/80 dark:border-secondary-dark dark:bg-secondary-dark dark:text-on-secondary-dark/80"
	case Info:
		return "border border-info bg-info text-on-info/80"
	case Success:
		return "border border-success bg-success text-on-success/80"
	case Warning:
		return "border border-warning bg-warning text-on-warning/80"
	case Danger:
		return "border border-danger bg-danger text-on-danger/80"
	default:
		return "border border-outline bg-surface-alt text-on-surface/80 dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark/80"
	}
}

// BorderClasses returns border classes if border is enabled
func (cfg Config) BorderClasses() string {
	if !cfg.Border {
		return ""
	}

	color := cfg.BorderColor
	if color == "" {
		switch cfg.Variant {
		case Info:
			color = "border-info"
		case Success:
			color = "border-success"
		case Warning:
			color = "border-warning"
		case Danger:
			color = "border-danger"
		default:
			color = "border-primary"
		}
	}

	return "border-2 " + color + " p-0.5"
}

// StatusClasses returns status dot classes
func (cfg Config) StatusClasses() string {
	switch cfg.Status {
	case StatusOffline:
		return "bg-outline dark:bg-outline-dark"
	case StatusInfo:
		return "bg-info"
	case StatusSuccess:
		return "bg-success"
	case StatusWarning:
		return "bg-warning"
	case StatusDanger:
		return "bg-danger"
	default:
		return ""
	}
}

// StatusSizeClasses returns status dot size based on avatar size
func (cfg Config) StatusSizeClasses() string {
	switch cfg.Size {
	case SizeXS:
		return "size-2"
	case SizeSM:
		return "size-2.5"
	case SizeLG:
		return "size-5"
	case SizeXL:
		return "size-6"
	default:
		return "size-4"
	}
}

// HasImage returns true if avatar uses an image
func (cfg Config) HasImage() bool {
	return cfg.Src != ""
}

// HasInitials returns true if avatar uses initials
func (cfg Config) HasInitials() bool {
	return cfg.Initials != ""
}

// GetInitials derives 1-2 character initials from a name, falling back to email.
// Examples: "John Doe" → "JD", "Alice" → "A", "" with "alice@x.com" → "a", "" with "" → "U"
func GetInitials(name, email string) string {
	if name != "" {
		parts := splitWords(name)
		if len(parts) >= 2 {
			return string(parts[0][0]) + string(parts[len(parts)-1][0])
		}
		if len(parts) == 1 && len(parts[0]) > 0 {
			return string(parts[0][0])
		}
	}
	if len(email) > 0 {
		return string(email[0])
	}
	return "U"
}

func splitWords(s string) []string {
	var result []string
	var current string
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
