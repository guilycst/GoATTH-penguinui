package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
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

func main() {
	mux := http.NewServeMux()

	// Serve the main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		Page().Render(r.Context(), w)
	})

	// HTMX endpoint: returns filtered/sorted table rows
	mux.HandleFunc("/api/breeds", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		q := r.URL.Query()
		search := strings.ToLower(q.Get("search"))
		group := q.Get("group")
		orderBy := q.Get("order_by")
		orderDir := q.Get("order_dir")

		// Filter
		var filtered []Dog
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

		// Sort
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

		// Render rows
		for _, d := range filtered {
			fmt.Fprintf(w, `<tr class="border-t border-gray-200">
				<td class="px-4 py-3 font-medium">%s</td>
				<td class="px-4 py-3">%s</td>
				<td class="px-4 py-3">%s</td>
				<td class="px-4 py-3">%s</td>
				<td class="px-4 py-3">%s</td>
			</tr>`, d.Breed, d.Group, d.Origin, d.Size, d.Temperament)
		}

		if len(filtered) == 0 {
			fmt.Fprint(w, `<tr><td colspan="5" class="px-4 py-8 text-center text-gray-400">No breeds found</td></tr>`)
		}
	})

	addr := ":3000"
	log.Printf("Dog breeds app running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
