package edgeos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// getHTTP creates http requests to download data
func getHTTP(o *object) *object {
	var (
		body []byte
		err  error
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest(o.Method, o.url, nil)
	if err != nil {
		o.r = strings.NewReader(fmt.Sprintf("Unable to form request for %s...", o.url))
		o.err = err
		return o
	}

	req.Header.Set("User-Agent", agent)
	resp, err = (&http.Client{}).Do(req)
	if err != nil {
		o.r = strings.NewReader(fmt.Sprintf("Unable to get response for %s...", o.url))
		o.err = err
		return o
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if len(body) == 0 {
		o.r = strings.NewReader(fmt.Sprintf("No data returned for %s...", o.url))
		o.err = err
		return o
	}

	o.r = bytes.NewBuffer(body)
	o.err = err

	return o
}
