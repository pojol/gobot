package bot

import (
	"net/http"
	"sync"
	"time"

	script "github.com/pojol/gobot/script/module"
	lua "github.com/yuin/gopher-lua"
)

type lStatePool struct {
	m     sync.Mutex
	saved []*botState
}

type botState struct {
	L *lua.LState

	httpMod   *script.HttpModule
	protoMod  *script.ProtoModule
	utilsMod  *script.UtilsModule
	base64Mod *script.Base64Module
	mgoMod    *script.MgoModule
}

func (pl *lStatePool) Get() *botState {
	pl.m.Lock()
	defer pl.m.Unlock()

	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]

	return x
}

func (pl *lStatePool) New() *botState {

	b := &botState{
		L:         lua.NewState(),
		httpMod:   script.NewHttpModule(&http.Client{Timeout: time.Second * 120}),
		protoMod:  &script.ProtoModule{},
		utilsMod:  &script.UtilsModule{},
		base64Mod: &script.Base64Module{},
		mgoMod:    &script.MgoModule{},
	}

	b.L.PreloadModule("proto", b.protoMod.Loader)
	b.L.PreloadModule("http", b.httpMod.Loader)
	b.L.PreloadModule("utils", b.utilsMod.Loader)
	b.L.PreloadModule("base64", b.base64Mod.Loader)
	b.L.PreloadModule("mgo", b.mgoMod.Loader)

	return b
}

func (pl *lStatePool) Put(bs *botState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, bs)
}

func (pl *lStatePool) Shutdown() {
	for _, bs := range pl.saved {
		bs.L.Close()
	}
}

// Global LState pool
var luaPool = &lStatePool{
	saved: make([]*botState, 0, 64),
}
