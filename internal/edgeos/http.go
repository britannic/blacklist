package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// getHTTP creates http requests to download data
func getHTTP(method, URL string) (io.Reader, error) {
	var (
		body []byte
		err  error
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest(method, URL, nil)
	if err != nil {
		return strings.NewReader(fmt.Sprintf("Unable to form request for %s...", URL)), err
	}

	req.Header.Set("User-Agent", agent)
	resp, err = (&http.Client{}).Do(req)
	if err != nil {
		return strings.NewReader(fmt.Sprintf("Unable to get response for %s...", URL)), err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if len(body) == 0 {
		return strings.NewReader(fmt.Sprintf("No data returned for %s...", URL)), err
	}

	return bytes.NewBuffer(body), err
}
