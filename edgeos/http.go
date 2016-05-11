package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPer implements the Do method to abstract HTTP requests
// type HTTPer interface {
// 	Do(*http.Request) (*http.Response, error)
// }

// HTTP implements the HTTPer interface
// type HTTP struct {
// 	HTTPer
// }

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

	// fmt.Println(string(body[:]))
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
