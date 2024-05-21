package utils

import (
	"net/http"
	"strings"

	"github.com/BetOnz-Company/betonz-go/internal/db"
)

func ParseRequest(request *http.Request) db.HttpRequest {
	return db.HttpRequest{
		Url:     request.URL.String(),
		Method:  request.Method,
		Headers: parseHeaders(request),
	}
}

func parseHeaders(request *http.Request) map[string]string {
	headerRecord := make(map[string]string)

	for key, value := range request.Header {
		key = strings.ToLower(key)
		if key != "cookie" {
			headerRecord[key] = strings.Join(value, ", ")
		}
	}

	return headerRecord
}
