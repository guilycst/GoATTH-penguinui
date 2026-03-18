package main

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/guilycst/GoATTH-penguinui/components/table"
)

// Dog represents a dog breed with metadata
type Dog struct {
	Breed       string
	Group       string
	Origin      string
	Size        string
	Temperament string
}

var breeds = []Dog{
	{Breed: "Labrador Retriever", Group: "Sporting", Origin: "Canada", Size: "Large", Temperament: "Friendly"},
	{Breed: "German Shepherd", Group: "Herding", Origin: "Germany", Size: "Large", Temperament: "Loyal"},
	{Breed: "Golden Retriever", Group: "Sporting", Origin: "Scotland", Size: "Large", Temperament: "Gentle"},
	{Breed: "French Bulldog", Group: "Non-Sporting", Origin: "France", Size: "Small", Temperament: "Playful"},
	{Breed: "Bulldog", Group: "Non-Sporting", Origin: "England", Size: "Medium", Temperament: "Calm"},
	{Breed: "Poodle", Group: "Non-Sporting", Origin: "Germany", Size: "Medium", Temperament: "Intelligent"},
	{Breed: "Beagle", Group: "Hound", Origin: "England", Size: "Small", Temperament: "Curious"},
	{Breed: "Rottweiler", Group: "Working", Origin: "Germany", Size: "Large", Temperament: "Confident"},
	{Breed: "Dachshund", Group: "Hound", Origin: "Germany", Size: "Small", Temperament: "Clever"},
	{Breed: "Yorkshire Terrier", Group: "Toy", Origin: "England", Size: "Small", Temperament: "Spirited"},
	{Breed: "Boxer", Group: "Working", Origin: "Germany", Size: "Large", Temperament: "Energetic"},
	{Breed: "Siberian Husky", Group: "Working", Origin: "Russia", Size: "Medium", Temperament: "Outgoing"},
	{Breed: "Shih Tzu", Group: "Toy", Origin: "China", Size: "Small", Temperament: "Affectionate"},
	{Breed: "Border Collie", Group: "Herding", Origin: "Scotland", Size: "Medium", Temperament: "Smart"},
	{Breed: "Doberman", Group: "Working", Origin: "Germany", Size: "Large", Temperament: "Alert"},
	{Breed: "Corgi", Group: "Herding", Origin: "Wales", Size: "Small", Temperament: "Happy"},
	{Breed: "Australian Shepherd", Group: "Herding", Origin: "USA", Size: "Medium", Temperament: "Active"},
	{Breed: "Cavalier King Charles", Group: "Toy", Origin: "England", Size: "Small", Temperament: "Graceful"},
	{Breed: "Great Dane", Group: "Working", Origin: "Germany", Size: "Large", Temperament: "Patient"},
	{Breed: "Chihuahua", Group: "Toy", Origin: "Mexico", Size: "Small", Temperament: "Charming"},
}

// columns defines the table headers — sortable columns get click-to-sort via HTMX
func columns() []table.Column {
	return []table.Column{
		{Key: "breed", Label: "Breed", Sortable: true},
		{Key: "group", Label: "Group", Sortable: true},
		{Key: "origin", Label: "Origin", Sortable: true},
		{Key: "size", Label: "Size", Sortable: true},
		{Key: "temperament", Label: "Temperament"},
	}
}

// dogsToRows converts a slice of Dogs into table.Row for the component
func dogsToRows(dogs []Dog) []table.Row {
	rows := make([]table.Row, len(dogs))
	for i, d := range dogs {
		rows[i] = table.Row{
			ID: d.Breed,
			Cells: map[string]table.Cell{
				"breed":       {Text: d.Breed},
				"group":       {Text: d.Group},
				"origin":      {Text: d.Origin},
				"size":        {Text: d.Size},
				"temperament": {Text: d.Temperament},
			},
		}
	}
	return rows
}

// filterAndSort applies search, group filter, and sort to the breed list
func filterAndSort(search, group, orderBy, orderDir string) []Dog {
	var filtered []Dog
	search = strings.ToLower(search)
	for _, d := range breeds {
		if group != "" && d.Group != group {
			continue
		}
		if search != "" &&
			!strings.Contains(strings.ToLower(d.Breed), search) &&
			!strings.Contains(strings.ToLower(d.Origin), search) &&
			!strings.Contains(strings.ToLower(d.Temperament), search) {
			continue
		}
		filtered = append(filtered, d)
	}
	if orderBy != "" {
		sort.SliceStable(filtered, func(i, j int) bool {
			var a, b string
			switch orderBy {
			case "breed":
				a, b = filtered[i].Breed, filtered[j].Breed
			case "group":
				a, b = filtered[i].Group, filtered[j].Group
			case "origin":
				a, b = filtered[i].Origin, filtered[j].Origin
			case "size":
				a, b = filtered[i].Size, filtered[j].Size
			}
			if orderDir == "desc" {
				return a > b
			}
			return a < b
		})
	}
	return filtered
}

func main() {
	mux := http.NewServeMux()

	// Serve the main page with the full table
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		dogs := filterAndSort("", "", "breed", "asc")
		Page(columns(), dogsToRows(dogs)).Render(r.Context(), w)
	})

	// HTMX endpoint: returns filtered/sorted table rows
	mux.HandleFunc("/api/breeds", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		q := r.URL.Query()
		dogs := filterAndSort(q.Get("search"), q.Get("group"), q.Get("order_by"), q.Get("order_dir"))
		cfg := table.Config{
			ID:           "breeds",
			Columns:      columns(),
			Rows:         dogsToRows(dogs),
			HTMXEndpoint: "/api/breeds",
			SortBy:       q.Get("order_by"),
			SortDir:      table.SortDir(q.Get("order_dir")),
		}
		for _, row := range cfg.Rows {
			table.TableRow(cfg, row).Render(r.Context(), w)
		}
	})

	addr := ":3000"
	log.Printf("Dog breeds app running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
