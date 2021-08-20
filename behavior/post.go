package behavior

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HTTPPost struct {
	URL    string
	Header map[string]string
}

func (p *HTTPPost) Do(in []byte) ([]byte, error) {

	var res *http.Response
	var out []byte

	client := http.Client{}

	req, err := http.NewRequest("POST", p.URL, bytes.NewBuffer(in))
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
		out, _ = ioutil.ReadAll(res.Body)
	} else {
		io.Copy(ioutil.Discard, res.Body)
		err = fmt.Errorf("http status %v url = %v err", res.Status, p.URL)
	}

ext:
	return out, err
}
