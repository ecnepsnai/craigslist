package craigslist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var _cookie string
var _bearerMap = map[int]string{}

func getCookie() (string, error) {
	if _cookie != "" {
		return _cookie, nil
	}

	req, err := http.NewRequest("HEAD", "https://api.craigslist.org/connection-check", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", eclUserAgent)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("http %d", response.StatusCode)
	}

	for _, c := range response.Cookies() {
		_cookie = fmt.Sprintf("%s=%s", c.Name, c.Value)
	}

	return _cookie, nil
}

func getBearer(areaID int) (string, error) {
	if bearer := _bearerMap[areaID]; bearer != "" {
		return bearer, nil
	}
	cookie, err := getCookie()
	if err != nil {
		return "", nil
	}

	body := &bytes.Buffer{}
	body.Write(providerCredHeaderValue)
	req, err := http.NewRequest("POST", "https://rapi.craigslist.org/v7/access-token", body)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"x-ecl-appversion": eclAppVersion,
		"Accept":           "application/json",
		"x-ecl-appname":    eclAppName,
		"x-ecl-doorkey":    eclDoorKey,
		"x-ecl-deviceid":   uuid.New().String(),
		"x-ecl-areaid":     fmt.Sprintf("%d", areaID),
		"Accept-Language":  "en-un",
		"x-ecl-devicename": randomString(16),
		"x-ecl-logid":      eclLogID,
		"Content-Length":   fmt.Sprintf("%d", body.Len()),
		"User-Agent":       eclUserAgent,
		"Content-Type":     "application/x-www-form-urlencoded",
		"Cookie":           cookie,
		"x-ecl-useragent":  eclUserAgent,
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("http %d", response.StatusCode)
	}

	type configResponseType struct {
		Data struct {
			Items []struct {
				AccessToken string `json:"accessToken"`
			} `json:"items"`
		} `json:"data"`
	}
	data := configResponseType{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", nil
	}

	bearer := data.Data.Items[0].AccessToken
	_bearerMap[areaID] = bearer
	return bearer, nil
}
