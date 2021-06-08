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
	Api    api.API
}

func (p *HTTPPost) Exec() error {

	var res *http.Response

	byt := p.Api.Marshal(p.Meta)

	client := http.Client{}

	req, err := http.NewRequest("POST", p.URL, bytes.NewBuffer(byt))
	if err != nil {
		goto ext
	}

	res, err = client.Do(req)
	if err != nil {
		goto ext
	}
	defer res.Body.Close()
	req.Body.Close()

	if res.StatusCode == http.StatusOK {

		body, _ := ioutil.ReadAll(res.Body)
		p.Api.Unmarshal(p.Meta, body, res.Header)

		err = p.Api.Assert(p.Meta)

	} else {
		io.Copy(ioutil.Discard, res.Body)
		err = fmt.Errorf("http status %v url = %v err", res.Status, p.URL)
	}

ext:
	return err
}
