package craigslist_test

import (
	"net/http"
	"testing"

	"github.com/ecnepsnai/craigslist"
)

func TestGetPostings(t *testing.T) {
	// Perform a search for 'vintage' in 'for sale - computers' in 'Vancouver BC'
	postings, err := craigslist.GetPostings("sya", "vintage", craigslist.LocationParams{
		AreaID:         16,
		Latitude:       49.2810,
		Longitude:      -123.0400,
		SearchDistance: 30,
	})
	if err != nil {
		t.Fatalf("Error getting postings: %s", err.Error())
	}

	if len(postings) == 0 {
		t.Errorf("No posts returned :(")
	}

	if !testing.Short() {
		for _, post := range postings {
			if len(post.Images) == 0 {
				continue
			}

			imageURL := post.ImageURLs()[0]
			resp, _ := http.Get(imageURL)
			if resp.StatusCode != 200 {
				t.Errorf("HTTP %d for image URL", resp.StatusCode)
			}
			break
		}
	}
}
