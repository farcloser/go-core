package network

import (
	"fmt"
	"net/http"
	"strings"
)

// Transport implements http.Transport with a RoundTrip that has baked-in defaults, notably for GitHub
// It is not meant to be instantiated directly, but rather obtained through Get().Transport()
type Transport struct {
	http.Transport
	TokenValue string
	TokenType  string
}

func (adt *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.TokenValue != "" {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", adt.TokenType, adt.TokenValue))
	}
	if strings.HasSuffix(req.Host, "github.com") {
		req.Header.Set("Accept", "application/json")
		// req.Header.Set("Content-Type", "application/json")
	}
	return adt.Transport.RoundTrip(req)
}
