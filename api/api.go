package api

import "net/http"

type API interface {
	Marshal(meta interface{}, param interface{}) []byte
	Unmarshal(meta interface{}, body []byte, header http.Header)
	Assert(meta interface{}) error
}
