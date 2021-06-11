package craigslist_test

import (
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

	didAnyReturnImages := false
	for _, r := range results {
		urls := r.ImageURLs()
		if len(urls) > 0 {
			didAnyReturnImages = true
			break
		}
	}
	if !didAnyReturnImages {
		// Note: this may TECHNICALLY fail if, by total chance, there's just a bunch of listings without pictures on
		// them.
		t.Errorf("No image URLs for any results")
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
