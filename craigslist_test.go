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
