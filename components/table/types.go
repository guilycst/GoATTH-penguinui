package table

import "github.com/a-h/templ"

// Variant represents table style variants
type Variant string

const (
	Default      Variant = "default"
	Striped      Variant = "striped"
	WithCheckbox Variant = "checkbox"
)

// SortDir represents the sort direction
type SortDir string

const (
	SortAsc  SortDir = "asc"
	SortDesc SortDir = "desc"
	SortNone SortDir = ""
)

// Column defines a table column header
type Column struct {
	// Key is the column identifier used to look up cell values in Row.Cells
	Key string
	// Label is the display text for the column header
	Label string
	// Sortable marks this column as sortable (renders clickable header)
	Sortable bool
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

// PaginationConfig holds pagination state
type PaginationConfig struct {
	// CurrentPage is the 1-indexed current page number
	CurrentPage int
	// TotalPages is the total number of pages
	TotalPages int
	// PerPage is the number of items per page
	PerPage int
}

// InfiniteScrollConfig holds infinite scroll state
type InfiniteScrollConfig struct {
	// NextPage is the next page number to load
	NextPage int
	// HasMore indicates if more rows are available
	HasMore bool
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

	// --- Sorting ---
	// SortBy is the currently sorted column key
	SortBy string
	// SortDir is the current sort direction ("asc" or "desc")
	SortDir SortDir

	// --- HTMX Integration ---
	// HTMXEndpoint is the base URL for HTMX requests (sorting, pagination, lazy load)
	HTMXEndpoint string
	// HTMXTarget overrides the default HTMX swap target (defaults to tbody ID)
	HTMXTarget string

	// --- Lazy Loading ---
	// LazyLoad loads the table body via HTMX on page load
	LazyLoad bool

	// --- Pagination ---
	// Pagination enables traditional pagination below the table
	Pagination *PaginationConfig

	// --- Infinite Scroll ---
	// InfiniteScroll enables loading more rows on scroll
	InfiniteScroll *InfiniteScrollConfig
}

// GetID returns the table ID, defaulting to "table"
func (cfg Config) GetID() string {
	if cfg.ID != "" {
		return cfg.ID
	}
	return "table"
}

// TbodyID returns the ID for the tbody element
func (cfg Config) TbodyID() string {
	return cfg.GetID() + "-tbody"
}

// PaginationID returns the ID for the pagination container
func (cfg Config) PaginationID() string {
	return cfg.GetID() + "-pagination"
}

// PaginationBaseURL returns the base URL for pagination links with per_page and sort params
func (cfg Config) PaginationBaseURL() string {
	url := cfg.HTMXEndpoint
	sep := "?"
	if cfg.Pagination != nil && cfg.Pagination.PerPage > 0 {
		url += sep + "per_page=" + itoa(cfg.Pagination.PerPage)
		sep = "&"
	}
	if cfg.SortBy != "" {
		url += sep + "order_by=" + cfg.SortBy + "&order_dir=" + string(cfg.SortDir)
	}
	return url
}

// HasSortableColumns returns true if any column is sortable
func (cfg Config) HasSortableColumns() bool {
	for _, col := range cfg.Columns {
		if col.Sortable {
			return true
		}
	}
	return false
}

// IsSortedBy returns true if the table is currently sorted by the given column
func (cfg Config) IsSortedBy(key string) bool {
	return cfg.SortBy == key
}

// NextSortDir returns the next sort direction when clicking a column header
func (cfg Config) NextSortDir(key string) SortDir {
	if cfg.SortBy != key || cfg.SortDir == SortNone {
		return SortAsc
	}
	if cfg.SortDir == SortAsc {
		return SortDesc
	}
	return SortAsc
}

// SortURL builds the HTMX URL for sorting by a given column
func (cfg Config) SortURL(key string) string {
	dir := cfg.NextSortDir(key)
	url := cfg.HTMXEndpoint + "?order_by=" + key + "&order_dir=" + string(dir)
	if cfg.Pagination != nil {
		url += "&per_page=" + itoa(cfg.Pagination.PerPage)
	}
	return url
}

// PageURL builds the HTMX URL for a specific page
func (cfg Config) PageURL(page int) string {
	url := cfg.HTMXEndpoint + "?page=" + itoa(page)
	if cfg.Pagination != nil {
		url += "&per_page=" + itoa(cfg.Pagination.PerPage)
	}
	if cfg.SortBy != "" {
		url += "&order_by=" + cfg.SortBy + "&order_dir=" + string(cfg.SortDir)
	}
	return url
}

// NextPageURL builds the HTMX URL for infinite scroll
func (cfg Config) NextPageURL() string {
	if cfg.InfiniteScroll == nil {
		return ""
	}
	url := cfg.HTMXEndpoint + "?page=" + itoa(cfg.InfiniteScroll.NextPage)
	if cfg.SortBy != "" {
		url += "&order_by=" + cfg.SortBy + "&order_dir=" + string(cfg.SortDir)
	}
	return url
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

// HeaderCellClasses returns CSS classes for a non-sortable header cell
func (cfg Config) HeaderCellClasses() string {
	return "p-4"
}

// SortableHeaderClasses returns CSS classes for a sortable header cell
func (cfg Config) SortableHeaderClasses(key string) string {
	base := "p-4 cursor-pointer select-none hover:bg-surface dark:hover:bg-surface-dark transition-colors"
	if cfg.IsSortedBy(key) {
		base += " text-primary dark:text-primary-dark"
	}
	return base
}

// CheckboxClasses returns CSS classes for checkboxes
func (cfg Config) CheckboxClasses() string {
	return "before:content[''] peer relative size-4 appearance-none overflow-hidden rounded border border-outline bg-surface before:absolute before:inset-0 checked:border-primary checked:before:bg-primary focus:outline-2 focus:outline-offset-2 focus:outline-outline-strong checked:focus:outline-primary active:outline-offset-0 dark:border-outline-dark dark:bg-surface-dark-alt dark:checked:border-primary-dark dark:checked:before:bg-primary-dark dark:focus:outline-outline-dark-strong dark:checked:focus:outline-primary-dark"
}

// PaginationPages returns the page numbers to display
func (p *PaginationConfig) PaginationPages() []int {
	if p == nil || p.TotalPages <= 0 {
		return nil
	}
	pages := make([]int, 0, p.TotalPages)
	for i := 1; i <= p.TotalPages; i++ {
		pages = append(pages, i)
	}
	return pages
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

// itoa converts an int to string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	if neg {
		digits = append(digits, '-')
	}
	// reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
