package behavior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pojol/apibot/utils"
	lua "github.com/yuin/gopher-lua"
)

// from https://github.com/cjoudrey/gluahttp

type httpModule struct {
	do func(req *http.Request) (*http.Response, error)
}

type empty struct{}

func NewHttpModule(client *http.Client) *httpModule {
	return NewHttpModuleWithDo(client.Do)
}

func NewHttpModuleWithDo(do func(req *http.Request) (*http.Response, error)) *httpModule {
	return &httpModule{
		do: do,
	}
}

func (h *httpModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":     h.get,
		"post":    h.post,
		"put":     h.put,
		"request": h.request,
	})
	registerHttpResponseType(mod, L)
	L.Push(mod)
	return 1
}

func (h *httpModule) get(L *lua.LState) int {
	return h.doRequestAndPush(L, "get", L.ToString(1), L.ToTable(2))
}

func (h *httpModule) post(L *lua.LState) int {
	return h.doRequestAndPush(L, "post", L.ToString(1), L.ToTable(2))
}

func (h *httpModule) put(L *lua.LState) int {
	return h.doRequestAndPush(L, "put", L.ToString(1), L.ToTable(2))
}

func (h *httpModule) request(L *lua.LState) int {
	return h.doRequestAndPush(L, L.ToString(1), L.ToString(2), L.ToTable(3))
}

func (h *httpModule) doRequest(L *lua.LState, method string, url string, options *lua.LTable) (*lua.LUserData, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		return nil, err
	}

	if ctx := L.Context(); ctx != nil {
		req = req.WithContext(ctx)
	}

	if options != nil {
		if reqCookies, ok := options.RawGet(lua.LString("cookies")).(*lua.LTable); ok {
			reqCookies.ForEach(func(key lua.LValue, value lua.LValue) {
				req.AddCookie(&http.Cookie{Name: key.String(), Value: value.String()})
			})
		}

		body := options.RawGet(lua.LString("body"))

		switch reqBody := body.(type) {
		case *lua.LTable:
			m, err := utils.Table2Map(reqBody)
			if err != nil {
				fmt.Println("table 2 map err", err.Error())
				return nil, err
			}
			byt, err := json.Marshal(m)
			if err != nil {
				fmt.Println("ltable marshal err", err.Error())
				return nil, err
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(byt))
			req.Header.Set("Content-Type", "application/json")
		}

		reqTimeout := options.RawGet(lua.LString("timeout"))
		if reqTimeout != lua.LNil {
			duration := time.Duration(0)
			switch reqTimeout.(type) {
			case lua.LNumber:
				duration = time.Second * time.Duration(int(reqTimeout.(lua.LNumber)))
			case lua.LString:
				duration, err = time.ParseDuration(string(reqTimeout.(lua.LString)))
				if err != nil {
					return nil, err
				}
			}
			ctx, cancel := context.WithTimeout(req.Context(), duration)
			req = req.WithContext(ctx)
			defer cancel()
		}

		// Basic auth
		if reqAuth, ok := options.RawGet(lua.LString("auth")).(*lua.LTable); ok {
			user := reqAuth.RawGetString("user")
			pass := reqAuth.RawGetString("pass")
			if !lua.LVIsFalse(user) && !lua.LVIsFalse(pass) {
				req.SetBasicAuth(user.String(), pass.String())
			} else {
				return nil, fmt.Errorf("auth table must contain no nil user and pass fields")
			}
		}

		// Set these last. That way the code above doesn't overwrite them.
		if reqHeaders, ok := options.RawGet(lua.LString("headers")).(*lua.LTable); ok {
			reqHeaders.ForEach(func(key lua.LValue, value lua.LValue) {
				req.Header.Set(key.String(), value.String())
			})
		}
	}

	res, err := h.do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return newHttpResponse(res, &body, len(body), L), nil
}

func (h *httpModule) doRequestAndPush(L *lua.LState, method string, url string, options *lua.LTable) int {
	response, err := h.doRequest(L, method, url, options)

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(response)
	return 1
}

func toTable(v lua.LValue) *lua.LTable {
	if lv, ok := v.(*lua.LTable); ok {
		return lv
	}
	return nil
}
