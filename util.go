package craigslist

import (
	"net/http"
	"strings"
)

const (
	eclAppName   = "craigslist mobile app"
	eclDoorKey   = "let me use the dev code to log in" // wtf??
	eclUserAgent = "CLApp/1.12.0/iOS unknown"
	eclLogID     = "6e8cfd5"
)

func httpGet(baseURL string, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	url := urlParamsToURL(baseURL, queryParams)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return http.DefaultClient.Do(req)
}

func urlParamsToURL(base string, params map[string]string) string {
	url := base + "?"
	paramsArr := make([]string, len(params))
	i := 0
	for key, value := range params {
		paramsArr[i] = key + "=" + value
		i++
	}
	url += strings.Join(paramsArr, "&")

	return url
}
