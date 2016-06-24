package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// GetHTTP creates http requests to download data
func GetHTTP(method, URL string) (io.Reader, error) {
	var (
		body []byte
		err  error
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest(method, URL, nil)
	if err == nil {
		req.Header.Set("User-Agent", agent)
		// req.Header.Add("Content-Type", "application/json")
		// debug(httputil.DumpRequestOut(req, true))
		resp, err = (&http.Client{}).Do(req)
		// resp, err = (&h{}).Do(req)
	} else {
		return bytes.NewBuffer([]byte(fmt.Sprintf("Unable to form request for %s...", URL))), err
	}

	if err == nil {
		defer resp.Body.Close()
		// debug(httputil.DumpResponse(resp, true))
		body, err = ioutil.ReadAll(resp.Body)
	}

	if len(body) == 0 {
		return bytes.NewBuffer([]byte(fmt.Sprintf("No data returned for %s...", URL))), err
	}

	return bytes.NewBuffer(body), err
}

// func debug(data []byte, err error) {
// 	switch {
// 	case !dbg:
// 		return
// 	case err == nil:
// 		fmt.Printf("Debug: %s\n\n", data)
// 	default:
// 		fmt.Printf("Error: %s\n\n", err)
// 	}
// }
