package triplet

import "encoding/json"

// Entry represents a single key-value-effect row.
type Entry struct {
	// Key is the entry key.
	Key string `json:"key"`
	// Value is the entry value.
	Value string `json:"value"`
	// Effect is the selected effect/category.
	Effect string `json:"effect"`
}

// EffectOption represents a selectable effect in the dropdown.
type EffectOption struct {
	// Value is the form submission value.
	Value string
	// Display is the display text.
	Display string
}

// Config holds configuration for the Triplet component.
// Triplet renders dynamic 3-field rows (key + value + dropdown) powered by Alpine.js.
// Submits as name[key]=value:effect.
type Config struct {
	// ID is the element ID
	ID string
	// Name is the form field name prefix (submits as name[key]=value:effect)
	Name string
	// Entries is the initial list of entries
	Entries []Entry
	// EffectOptions is the list of dropdown options for the third column
	EffectOptions []EffectOption
	// DefaultEffect is the default value for new rows (defaults to first option's value)
	DefaultEffect string
	// KeyPlaceholder is shown in each key input
	KeyPlaceholder string
	// ValuePlaceholder is shown in each value input
	ValuePlaceholder string
	// AddLabel is the "Add row" button text (default: "Add row")
	AddLabel string
	// Disabled prevents adding/removing entries
	Disabled bool
	// Class allows additional CSS classes on the container
	Class string
}

// GetAddLabel returns the add button label with default
func (c Config) GetAddLabel() string {
	if c.AddLabel != "" {
		return c.AddLabel
	}
	return "Add row"
}

// GetDefaultEffect returns the default effect for new rows
func (c Config) GetDefaultEffect() string {
	if c.DefaultEffect != "" {
		return c.DefaultEffect
	}
	if len(c.EffectOptions) > 0 {
		return c.EffectOptions[0].Value
	}
	return ""
}

// AlpineData returns the x-data JSON string for Alpine.js initialization
func (c Config) AlpineData() string {
	entries := c.Entries
	if entries == nil {
		entries = []Entry{}
	}
	b, _ := json.Marshal(struct {
		Name    string  `json:"name"`
		Entries []Entry `json:"entries"`
	}{Name: c.Name, Entries: entries})
	return string(b)
}

// ContainerClasses returns CSS classes for the outer container
func (c Config) ContainerClasses() string {
	base := "flex flex-col gap-2"
	if c.Class != "" {
		return base + " " + c.Class
	}
	return base
}
