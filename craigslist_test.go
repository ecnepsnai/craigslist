package craigslist_test

import (
	"net/http"
	"testing"

	"github.com/ecnepsnai/craigslist"
)

func TestSearchAndGet(t *testing.T) {
	// Perform a search for 'lenovo thinkpad' in 'for sale - computers' in 'Vancouver BC'
	results, err := craigslist.Search("sya", "lenovo thinkpad", craigslist.LocationParams{
		AreaID:         16,
		Latitude:       49.2810,
		Longitude:      -123.0400,
		SearchDistance: 30,
	})
	if err != nil {
		t.Fatalf("Error getting results: %s", err.Error())
	}

	if len(results) == 0 {
		t.Fatalf("Search should return results but did not")
	}

	var imageURL string
	for _, r := range results {
		urls := r.ImageURLs()
		if len(urls) > 0 {
			imageURL = urls[0]
			break
		}
	}
	if imageURL == "" {
		// Note: this may TECHNICALLY fail if, by total chance, there's just a bunch of listings without pictures on
		// them.
		t.Fatalf("No image URLs for any results")
	}
	resp, err := http.Get(imageURL)
	if err != nil {
		t.Fatalf("Error getting image URL: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Error getting image URL: http %d", resp.StatusCode)
	}

	// Get the details of the first result
	posting, err := results[0].Posting()
	if err != nil {
		t.Fatalf("Error getting posting: %s", err.Error())
	}
	if posting == nil {
		t.Fatalf("Should return posting")
	}

	t.Logf("Would you like to buy a '%s' for $%d?", posting.Title, posting.Price)
}
