package craigslist_test

import (
	"fmt"

	"github.com/ecnepsnai/craigslist"
)

func ExampleSearch() {
	// Perform a service for 'vintage' in the computers for-sale category 'sys' in Vancouver, BC
	results, err := craigslist.Search("sya", "vintage", craigslist.LocationParams{
		AreaID:         16,
		Latitude:       49.2810,
		Longitude:      -123.0400,
		SearchDistance: 30,
	})
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		fmt.Printf("No results!")
	}

	// Results only contains summary information of a posting. Call 'Posting()' on a result
	// to get the full details of a result
	posting, err := results[0].Posting()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s - $%d\n", posting.Title, posting.Price)
	fmt.Printf("--------\n")
	fmt.Printf("%s\n", posting.Body)
}
