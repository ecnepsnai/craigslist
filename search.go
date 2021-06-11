package craigslist

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Result describes a single craigslist post
type Result struct {
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

// ImageURLs get an array of URLs that contain all images in this posting at or around 600x450
func (r Result) ImageURLs() []string {
	urls := make([]string, len(r.Images))
	for i, id := range r.Images {
		path := strings.Split(id, ":")[1]
		urls[i] = fmt.Sprintf("https://images.craigslist.org/%s_600x450.jpg", path)
	}

	return urls
}

// Posting get the full posting for the result
func (r Result) Posting() (*Posting, error) {
	return GetPosting(r.PostingID, r.CategoryAbbr, r.Location)
}

// PostingLocation describes the location for a craigslist post
type PostingLocation struct {
	AreaID      int     `json:"areaId"`
	Hostname    string  `json:"hostname"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	SubareaAbbr string  `json:"subareaAbbr"`
}

// LocationParams describes parameters for specifying the location for a craigslist query
type LocationParams struct {
	AreaID         int
	Latitude       float32
	Longitude      float32
	SearchDistance int
}

// Search will return an array of all postings for the given category. A query can be included to search for
// postings matching a specific query.
//
// The category must be the 3-letter category code. See https://github.com/ecnepsnai/craigslist/blob/main/categories.md
// for a list of all possible categories.
// The query can be an empty string to list all recent posts for the category.
// The location must is used to refine the craigslist search to your specific area, and is required
//
// Returns a slice of results, which are not as detailed as a full posting.
func Search(category string, query string, location LocationParams) ([]Result, error) {
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

	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}
	bearer, err := getBearer(location.AreaID)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Accept":                   "application/json",
		"x-ecl-devicecountry":      "US",
		"x-ecl-devicemanufacturer": "Apple",
		"x-ecl-devicemodel":        "iPad Pro 12.9-inch (3rd generation)",
		"x-ecl-systemname":         "iOS",
		"x-ecl-appversion":         eclAppVersion,
		"x-ecl-devicelocale":       "en",
		"x-ecl-systemversion":      "14.6",
		"x-ecl-devicebrand":        "Apple",
		"x-ecl-appname":            eclAppName,
		"x-ecl-doorkey":            eclDoorKey,
		"x-ecl-deviceid":           uuid.New().String(),
		"x-ecl-areaid":             fmt.Sprintf("%d", location.AreaID),
		"x-ecl-logid":              eclLogID,
		"x-ecl-useragent":          eclUserAgent,
		"x-ecl-devicename":         randomString(16),
		"User-Agent":               eclUserAgent,
		"Cookie":                   cookie,
		"Authorization":            "Bearer " + bearer,
	}

	response, err := httpGet("https://sapi.craigslist.org/v7/postings/"+category+"/search", queryParams, headers)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http %d", response.StatusCode)
	}

	type postingsResponseType struct {
		Data struct {
			Items []Result `json:"items"`
		} `json:"data"`
	}

	reply := postingsResponseType{}
	if err := json.NewDecoder(response.Body).Decode(&reply); err != nil {
		return nil, err
	}

	return reply.Data.Items, nil
}
