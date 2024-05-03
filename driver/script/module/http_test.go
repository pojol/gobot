package script

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHttpPut(t *testing.T) {
	req, err := http.NewRequest(strings.ToUpper("put"), "http://localhost:8500/v1/kv/apilist", nil)
	if err != nil {
		panic(err)
	}

	req.Body = ioutil.NopCloser(bytes.NewBufferString(`[{"url":"/login/user.pwdLogin","req":{"Acc":"bot001"},"res":{"UCharId":"5783100001","Token":"619cb9a0098cd500019a5f94"},"desc":"账号密码登陆"}]`))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println("status", res.Status)
	fmt.Println(ioutil.ReadAll(res.Body))
}
