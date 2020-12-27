package craigslist

import (
	"encoding/json"
	"fmt"

	"github.com/ecnepsnai/security"
	"github.com/google/uuid"
)

// Posting describes a single craigslist post
type Posting struct {
	CategoryAbbr string          `json:"categoryAbbr"`
	CategoryID   string          `json:"categoryId"`
	DedupeKey    string          `json:"dedupeKey"`
	Images       []string        `json:"images"`
	Location     PostingLocation `json:"location"`
	PostedDate   int             `json:"postedDate"`
	PostingID    int             `json:"postingId"`
	Price        int             `json:"price"`
	Title        string          `json:"title"`
}

// ImageURLs get an array of URLs that contain all images in this posting
func (p Posting) ImageURLs() []string {
	urls := make([]string, len(p.Images))
	for i, id := range p.Images {
		urls[i] = fmt.Sprintf("https://images.craigslist.org/%s_600x450.jpg", id)
	}

	return urls
}

// PostingLocation describes the location for a craigslist post
type PostingLocation struct {
	AreaID      int     `json:"areaId"`
	Hostname    string  `json:"hostname"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	SubareaAbbr string  `json:"subareaAbbr"`
}

type postingsResponseType struct {
	Data struct {
		Items []Posting `json:"items"`
	} `json:"data"`
}

// LocationParams describes parameters for specifying the location for a craigslist query
type LocationParams struct {
	AreaID         int
	Latitude       float32
	Longitude      float32
	SearchDistance int
}

// GetPostings will return an array of all postings for the given category. A query can be included to search for
// postings matching a specific query.
//
// The category must be the 3-letter category code. See https://github.com/ecnepsnai/craigslist/blob/main/categories.md
// for a list of all possible categories.
// The query can be an empty string to list all recent posts for the category.
// The location must is used to refine the craigslist search to your specific area, and is required
func GetPostings(category string, query string, location LocationParams) ([]Posting, error) {
	queryParams := map[string]string{
		"area_id":         fmt.Sprintf("%d", location.AreaID),
		"batchSize":       "100",
		"lat":             fmt.Sprintf("%f", location.Latitude),
		"lon":             fmt.Sprintf("%f", location.Longitude),
		"search_distance": "30",
		"startIndex":      "0",
		"lang":            "en",
		"cc":              "us",
	}

	if query != "" {
		queryParams["query"] = query
	}

	headers := map[string]string{
		"Accept":           "application/json",
		"x-ecl-appname":    eclAppName,
		"x-ecl-doorkey":    eclDoorKey,
		"x-ecl-deviceid":   uuid.New().String(),
		"x-ecl-areaid":     fmt.Sprintf("%d", location.AreaID),
		"x-ecl-logid":      eclLogID,
		"x-ecl-useragent":  eclUserAgent,
		"x-ecl-devicename": security.RandomString(16),
	}

	response, err := httpGet("https://sapi.craigslist.org/v5/postings/"+category+"/search", queryParams, headers)
	if err != nil {
		return nil, err
	}

	reply := postingsResponseType{}
	if err := json.NewDecoder(response.Body).Decode(&reply); err != nil {
		return nil, err
	}

	return reply.Data.Items, nil
}
