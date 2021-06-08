package api

import "net/http"

type API interface {
	Marshal() []byte
	Unmarshal([]byte, http.Header)
	Assert() error
}
