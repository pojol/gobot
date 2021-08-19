package behavior

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pojol/apibot/api"
)

type HTTPPost struct {
	URL    string
	Header map[string]string
	Meta   interface{}
	Param  interface{}
	Api    interface{}
}

func (p *HTTPPost) Do() error {

	var res *http.Response

	api := p.Api.(api.API)
	byt := api.Marshal(p.Meta, p.Param)

	client := http.Client{}

	req, err := http.NewRequest("POST", p.URL, bytes.NewBuffer(byt))
	if err != nil {
		fmt.Println("http.request", err.Error())
		goto ext
	}

	res, err = client.Do(req)
	if err != nil {
		fmt.Println("client.Do", err.Error())
		goto ext
	}
	defer res.Body.Close()
	req.Body.Close()

	if res.StatusCode == http.StatusOK {

		body, _ := ioutil.ReadAll(res.Body)
		api.Unmarshal(p.Meta, body, res.Header)

		//err = api.Assert(p.Meta)

	} else {
		io.Copy(ioutil.Discard, res.Body)
		err = fmt.Errorf("http status %v url = %v err", res.Status, p.URL)
	}

ext:
	return err
}
