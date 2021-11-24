package script

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestBase64Encode(t *testing.T) {
	base64Mod := Base64Module{}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", base64Mod.Loader)

	L.DoString(`
		local base64 = require("base64")
		local val = "W3sidXJsIjoiL2xvZ2luL3VzZXIucHdkTG9naW4iLCJyZXEiOnsiQWNjIjoiYm90MDAxIn0sInJlcyI6eyJVQ2hhcklkIjoiNTc4MzEwMDAwMSIsIlRva2VuIjoiNjE5ZGFjY2FkMmY3YzcwMDAxZmY1OWM1In0sImRlc2MiOiLotKblj7flr4bnoIHnmbvpmYYifV0="

		s, err = base64.decode(val)
		print("base64 decode", s, err)

		if base64.encode(s) ~= val then
			error("base64 test fail")
		end
		
	`)
}
