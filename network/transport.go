package network

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Transport implements http.Transport with a RoundTrip that has baked-in defaults, notably for GitHub
// It is not meant to be instantiated directly, but rather obtained through Get().Transport().
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
		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	resp, err := adt.Transport.RoundTrip(req)
	if err != nil {
		err = errors.Join(ErrRoundTrip, err)
	}

	return resp, err
}
