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

	if req, err = http.NewRequest(o.Method, o.url, nil); err != nil {
		o.error(fmt.Sprintf("Unable to form request for %s. Error: %v", o.url, err))
		o.r, o.err = strings.NewReader(fmt.Sprintf("Unable to form request for %s...", o.url)), err
		return o
	}

	o.Log.Info(fmt.Sprintf("Downloading %s source %s", o.area(), o.name))

	req.Header.Set("User-Agent", agent)
	if resp, err = (&http.Client{}).Do(req); err != nil {
		o.error(fmt.Sprintf("Unable to get response for %s. Error: %v", o.url, err))
		o.r, o.err = strings.NewReader(fmt.Sprintf("Unable to get response for %s...", o.url)), err
		return o
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if len(body) == 0 {
		o.warning(fmt.Sprintf("No data returned for %s. Error: %v", o.url, err))
		o.r, o.err = strings.NewReader(fmt.Sprintf("No data returned for %s...", o.url)), err
		return o
	}

	o.r, o.err = bytes.NewBuffer(body), err

	return o
}
