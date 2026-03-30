package tagslist

import "encoding/json"

// Config holds configuration for the TagsList component.
// TagsList renders a dynamic list of text inputs powered by Alpine.js.
// Users can add and remove string values. Form submission uses indexed names: name[0], name[1], ...
type Config struct {
	// ID is the element ID
	ID string
	// Name is the form field name prefix (submits as name[0], name[1], ...)
	Name string
	// Values is the initial list of tag values
	Values []string
	// Placeholder is shown in each input (e.g. "e.g. prod, critical")
	Placeholder string
	// AddLabel is the "Add tag" button text (default: "Add tag")
	AddLabel string
	// Disabled prevents adding/removing tags
	Disabled bool
	// Class allows additional CSS classes on the container
	Class string
}

// GetAddLabel returns the add button label with default
func (c Config) GetAddLabel() string {
	if c.AddLabel != "" {
		return c.AddLabel
	}
	return "Add tag"
}

// AlpineData returns the x-data JSON string for Alpine.js initialization
func (c Config) AlpineData() string {
	values := c.Values
	if values == nil {
		values = []string{}
	}
	b, _ := json.Marshal(struct {
		Items []string `json:"items"`
		Name  string   `json:"name"`
	}{Items: values, Name: c.Name})
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
