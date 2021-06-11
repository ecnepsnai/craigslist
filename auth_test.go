package craigslist

import "testing"

func TestAuth(t *testing.T) {
	bearer, err := getBearer(16)
	if err != nil {
		t.Fatalf("Error getting bearer: %s", err.Error())
	}
	if bearer == "" {
		t.Fatalf("No bearer returned")
	}
}
