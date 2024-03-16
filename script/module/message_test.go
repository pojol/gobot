package script

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	RegisterMessageType(L)

	if err := L.DoString(`
        
		

		body = proto.unmarshal("LoginGuestRes", msgbody)
		print(msgId, body)

    `); err != nil {
		panic(err)
	}
}
