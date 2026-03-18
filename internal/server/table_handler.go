package server

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/guilycst/GoATTH-penguinui/components/table"
)

// tableRecord is the server-side data model for demo table rows
type tableRecord struct {
	ID         string
	Name       string
	Email      string
	Membership string
}

// allRecords returns the full dataset for table demos
func allRecords() []tableRecord {
	return []tableRecord{
		{ID: "2335", Name: "Alice Brown", Email: "alice.brown@penguinui.com", Membership: "Silver"},
		{ID: "2338", Name: "Bob Johnson", Email: "johnson.bob@penguinui.com", Membership: "Gold"},
		{ID: "2342", Name: "Sarah Adams", Email: "s.adams@penguinui.com", Membership: "Gold"},
		{ID: "2345", Name: "Alex Martinez", Email: "alex.martinez@penguinui.com", Membership: "Gold"},
		{ID: "2346", Name: "Ryan Thompson", Email: "ryan.thompson@penguinui.com", Membership: "Silver"},
		{ID: "2349", Name: "Emily Rodriguez", Email: "emily.rodriguez@penguinui.com", Membership: "Gold"},
		{ID: "2350", Name: "James Wilson", Email: "james.wilson@penguinui.com", Membership: "Silver"},
		{ID: "2351", Name: "Sophia Chen", Email: "sophia.chen@penguinui.com", Membership: "Gold"},
		{ID: "2352", Name: "Michael Davis", Email: "m.davis@penguinui.com", Membership: "Silver"},
		{ID: "2353", Name: "Olivia Taylor", Email: "olivia.taylor@penguinui.com", Membership: "Gold"},
		{ID: "2354", Name: "Daniel Lee", Email: "daniel.lee@penguinui.com", Membership: "Silver"},
		{ID: "2355", Name: "Emma Harris", Email: "emma.harris@penguinui.com", Membership: "Gold"},
	}
}

func recordToRow(rec tableRecord) table.Row {
	return table.Row{
		ID: rec.ID,
		Cells: map[string]table.Cell{
			"id":         {Text: rec.ID},
			"name":       {Text: rec.Name},
			"email":      {Text: rec.Email},
			"membership": {Text: rec.Membership},
		},
	}
}

func recordsToRows(recs []tableRecord) []table.Row {
	rows := make([]table.Row, len(recs))
	for i, rec := range recs {
		rows[i] = recordToRow(rec)
	}
	return rows
}

func sortRecords(recs []tableRecord, orderBy string, orderDir string) {
	sort.SliceStable(recs, func(i, j int) bool {
		var a, b string
		switch orderBy {
		case "id":
			a, b = recs[i].ID, recs[j].ID
		case "name":
			a, b = strings.ToLower(recs[i].Name), strings.ToLower(recs[j].Name)
		case "email":
			a, b = strings.ToLower(recs[i].Email), strings.ToLower(recs[j].Email)
		case "membership":
			a, b = strings.ToLower(recs[i].Membership), strings.ToLower(recs[j].Membership)
		default:
			return false
		}
		if orderDir == "desc" {
			return a > b
		}
		return a < b
	})
}

func (s *Server) handleTableRows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	q := r.URL.Query()
	variant := q.Get("variant")
	orderBy := q.Get("order_by")
	orderDir := q.Get("order_dir")
	pageStr := q.Get("page")
	perPageStr := q.Get("per_page")

	records := allRecords()

	// Apply sorting
	if orderBy != "" {
		if orderDir == "" {
			orderDir = "asc"
		}
		sortRecords(records, orderBy, orderDir)
	}

	// Simulate server latency for lazy load / infinite scroll demos
	if variant == "lazy" || variant == "infinite" {
		time.Sleep(500 * time.Millisecond)
	}

	// Parse pagination params
	page := 1
	perPage := 3
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}

	totalPages := (len(records) + perPage - 1) / perPage

	// Paginate
	start := (page - 1) * perPage
	if start >= len(records) {
		start = 0
		page = 1
	}
	end := start + perPage
	if end > len(records) {
		end = len(records)
	}
	pageRecords := records[start:end]
	rows := recordsToRows(pageRecords)

	hasMore := end < len(records)
	nextPage := page + 1

	cfg := table.Config{
		Columns: []table.Column{
			{Key: "id", Label: "CustomerID"},
			{Key: "name", Label: "Name"},
			{Key: "email", Label: "Email"},
			{Key: "membership", Label: "Membership"},
		},
		Rows:    rows,
		SortBy:  orderBy,
		SortDir: table.SortDir(orderDir),
	}

	// For infinite scroll, render rows without tbody wrapper (appended to existing tbody)
	if variant == "infinite" {
		if hasMore {
			cfg.HTMXEndpoint = "/api/components/table/rows?variant=infinite"
			cfg.InfiniteScroll = &table.InfiniteScrollConfig{
				NextPage: nextPage,
				HasMore:  true,
			}
		}
		table.TableRows(cfg).Render(r.Context(), w)
		return
	}

	// For pagination, render rows as tbody inner HTML
	if pageStr != "" || variant == "" {
		cfg.Pagination = &table.PaginationConfig{
			CurrentPage: page,
			TotalPages:  totalPages,
			PerPage:     perPage,
		}
	}

	// Render just the table rows (tbody inner content)
	for _, row := range rows {
		table.TableRow(cfg, row).Render(r.Context(), w)
	}
}
