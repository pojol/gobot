package script

import "net/http"

type IScriptObject interface {
	Marshal() []byte
	Unmarshal([]byte, http.Header)
	Assert() error
}
