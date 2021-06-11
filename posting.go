package craigslist

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Posting describes the structure for a single post
type Posting struct {
	PostedDate   int64              `json:"postedDate"`
	SeeMyOther   int                `json:"seeMyOther"`
	URL          string             `json:"url"`
	Location     PostingLocation    `json:"location"`
	CategoryID   int                `json:"categoryId"`
	UpdatedDate  int64              `json:"updatedDate"`
	Attributes   []PostingAttribute `json:"attributes"`
	Category     string             `json:"category"`
	Body         string             `json:"body"`
	PostingID    int64              `json:"postingId"`
	Section      string             `json:"section"`
	CategoryAbbr string             `json:"categoryAbbr"`
	Title        string             `json:"title"`
	Images       []string           `json:"images"`
	Price        int                `json:"price"`
}

// PostingAttribute describes the structure for a single attribute on a posting
type PostingAttribute struct {
	SpecialType string `json:"specialType"`
	Label       string `json:"label"`
	Value       string `json:"value"`
}

// GetPosting get the full posting details for the given post.
// It's recommended that you don't call this method directly, but rather call 'Posting()' on
// a `craigslist.Result` object instead.
func GetPosting(id int, categoryAbbr string, location PostingLocation) (*Posting, error) {
	queryParams := map[string]string{
		"lang": "en",
		"cc":   "us",
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

	response, err := httpGet("https://api.craigslist.org/v7/postings/"+location.Hostname+"/"+location.SubareaAbbr+"/"+categoryAbbr+"/"+fmt.Sprintf("%d", id), queryParams, headers)
	if err != nil {
		return nil, err
	}

	type postingResponseType struct {
		Data struct {
			Items []Posting `json:"items"`
		} `json:"data"`
	}

	reply := postingResponseType{}
	if err := json.NewDecoder(response.Body).Decode(&reply); err != nil {
		return nil, err
	}

	if len(reply.Data.Items) == 0 {
		return nil, fmt.Errorf("no posting found")
	}

	post := reply.Data.Items[0]
	return &post, nil
}
