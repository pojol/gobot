package pool

import (
	"sync"

	script "github.com/pojol/gobot/script/module"
	lua "github.com/yuin/gopher-lua"
)

type lStatePool struct {
	m     sync.Mutex
	saved []*BotState
}

type BotState struct {
	L         *lua.LState
	HttpMod   *script.HttpModule
	TCPMod    *script.TCPModule
	protoMod  *script.ProtoModule
	utilsMod  *script.UtilsModule
	base64Mod *script.Base64Module
	mgoMod    *script.MgoModule
	md5Mod    *script.MD5Module
}

func (pl *lStatePool) Get() *BotState {
	pl.m.Lock()
	defer pl.m.Unlock()

	n := len(pl.saved)
	if n == 0 {
		return _new_state()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]

	return x
}

func _new_state() *BotState {

	b := &BotState{
		L:         lua.NewState(),
		HttpMod:   script.NewHttpModule(),
		TCPMod:    script.NewTCPModule(),
		protoMod:  &script.ProtoModule{},
		utilsMod:  &script.UtilsModule{},
		base64Mod: &script.Base64Module{},
		mgoMod:    &script.MgoModule{},
		md5Mod:    &script.MD5Module{},
	}

	b.L.PreloadModule("proto", b.protoMod.Loader)
	b.L.PreloadModule("http", b.HttpMod.Loader)
	b.L.PreloadModule("tcpconn", b.TCPMod.Loader)
	b.L.PreloadModule("utils", b.utilsMod.Loader)
	b.L.PreloadModule("base64", b.base64Mod.Loader)
	b.L.PreloadModule("mgo", b.mgoMod.Loader)
	b.L.PreloadModule("md5", b.md5Mod.Loader)

	return b
}

func (pl *lStatePool) Put(bs *BotState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, bs)
}

func (pl *lStatePool) Shutdown() {
	for _, bs := range pl.saved {
		bs.L.Close()
	}
}

func GetState() *BotState {
	return luaPool.Get()
}

func PutState(state *BotState) {
	luaPool.Put(state)
}

func NewState() *BotState {
	return _new_state()
}

func FreeState(state *BotState) {
	state.L.Close()
}

// Global LState pool
var luaPool = &lStatePool{
	saved: make([]*BotState, 0, 1024),
}
