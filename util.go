package craigslist

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
)

const (
	eclAppName    = "craigslist mobile app"
	eclDoorKey    = "let me use the dev code to log in" // wtf??
	eclUserAgent  = "CLApp/1.14.2/iOS unknown"
	eclAppVersion = "1.14.2-20210412-152100-94d307e1"
	eclLogID      = "2a1c34f"
)

func httpGet(baseURL string, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	clURL := urlParamsToURL(baseURL, queryParams)
	req, err := http.NewRequest("GET", clURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return http.DefaultClient.Do(req)
}

func urlParamsToURL(base string, params map[string]string) string {
	clURL := base + "?"
	paramsArr := make([]string, len(params))
	i := 0
	for key, value := range params {
		paramsArr[i] = key + "=" + url.PathEscape(value)
		i++
	}
	clURL += strings.Join(paramsArr, "&")
	return clURL
}

func randomString(length int) string {
	randB := make([]byte, length)
	rand.Read(randB)
	return hex.EncodeToString(randB)
}
