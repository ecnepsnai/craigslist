package craigslist

import (
	_ "embed"
)

// The API requires that you provide a URL-encoded certificate in a header.
// The certificate is the "Mac App Store and iTunes Store Receipt Signing" certificate.

//go:embed cert.txt
var providerCredHeaderValue []byte
