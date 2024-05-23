package azure

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"gorm.io/datatypes"
)

func realIP(r *http.Request) string {
	// Get the remote address from the request
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// If X-Forwarded-For header is set, it may contain a comma-separated list of IP addresses.
		// The client's IP address is usually the first one in the list.
		remoteAddr = ip
	}
	return remoteAddr
}

func toResponseHeaders(headers http.Header) *ResponseHeaders {
	respHeaders := &ResponseHeaders{}

	respHeaders.ContentLength = []string{headers.Get("Content-Length")}
	respHeaders.OperationLocation = []string{headers.Get("Operation-Location")}
	respHeaders.XEnvoyUpstreamServiceTime = []string{headers.Get("X-Envoy-Upstream-Service-Time")}
	respHeaders.APIMRequestID = []string{headers.Get("Apim-Request-Id")}
	respHeaders.StrictTransportSecurity = []string{headers.Get("Strict-Transport-Security")}
	respHeaders.XContentTypeOptions = []string{headers.Get("X-Content-Type-Options")}
	respHeaders.XMsRegion = []string{headers.Get("X-Ms-Region")}
	respHeaders.Date = []string{headers.Get("Date")}

	return respHeaders
}

func parseHeadersToMap(headers http.Header) (datatypes.JSON, error) {
	reqHeaders := make(map[string]interface{})

	for k, v := range headers {
		if len(v) > 0 {
			reqHeaders[k] = v
		}
	}

	// Convert slice of maps to JSON
	jsonData, err := json.Marshal(reqHeaders)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(jsonData), nil
}

func getAPIMRequestID(urlString string) string {
	// Parse the URL
	parsedURL, _ := url.Parse(urlString)

	// Extract the desired path segment
	return path.Base(parsedURL.Path)
}
