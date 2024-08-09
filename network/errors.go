package network

import "errors"

var (
	ErrRoundTrip                 = errors.New("round trip error")
	ErrInterfacesRetrievalFailed = errors.New("retrieving interfaces failed")
)
