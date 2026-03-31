package keyvalue

import "encoding/json"

// Entry represents a single key-value pair.
type Entry struct {
	// Key is the entry key.
	Key string `json:"key"`
	// Value is the entry value.
	Value string `json:"value"`
}

// Config holds configuration for the KeyValue component.
// KeyValue renders dynamic key-value pair rows powered by Alpine.js.
// Submits as name[key]=value.
type Config struct {
	// ID is the element ID
	ID string
	// Name is the form field name prefix (submits as name[key]=value)
	Name string
	// Entries is the initial list of key-value pairs
	Entries []Entry
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
