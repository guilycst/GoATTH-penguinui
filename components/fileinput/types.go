package fileinput

import "github.com/a-h/templ"

// Config holds configuration for the file input component
type Config struct {
	// ID is the HTML id for the file input element
	ID string
	// Name is the form field name
	Name string
	// Label text displayed above the drop zone (e.g. "Cover Picture")
	Label string
	// Accept restricts file types (e.g. "image/*", ".pdf,.doc")
	Accept string
	// HelperText shown below the drop zone (e.g. "PNG, JPG, WebP - Max 5MB")
	HelperText string
	// Required marks the input as required
	Required bool
	// Disabled disables the input
	Disabled bool
	// Class allows additional CSS classes on the outer container
	Class string
	// Attrs are extra attributes applied to the <input> element (e.g. hx-post, x-on:change)
	Attrs templ.Attributes
}

// ContainerClasses returns CSS classes for the outermost wrapper div
func (cfg Config) ContainerClasses() string {
	base := "flex w-full max-w-xl flex-col gap-1 text-center"
	if cfg.Class != "" {
		return base + " " + cfg.Class
	}
	return base
}

// LabelClasses returns CSS classes for the label text above the drop zone
func (cfg Config) LabelClasses() string {
	return "w-fit pl-0.5 text-sm text-on-surface dark:text-on-surface-dark"
}

// DropZoneClasses returns the static (non-dynamic) CSS classes for the drop zone
func (cfg Config) DropZoneClasses() string {
	base := "flex w-full flex-col items-center justify-center gap-2 rounded-radius border border-dashed p-8 text-on-surface dark:text-on-surface-dark"
	if cfg.Disabled {
		return base + " opacity-50 cursor-not-allowed"
	}
	return base
}

// BrowseLabelClasses returns CSS classes for the "Browse" label link
func (cfg Config) BrowseLabelClasses() string {
	return "font-medium text-primary group-focus-within:underline dark:text-primary-dark cursor-pointer"
}
