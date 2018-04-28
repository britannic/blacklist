package edgeos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// getHTTP creates http requests to download data
func getHTTP(s *source) *source {
	var (
		body []byte
		err  error
		resp *http.Response
		req  *http.Request
	)

	if req, err = http.NewRequest(s.Method, s.url, nil); err != nil {
		str := fmt.Sprintf("Unable to form request for %s", s.url)
		s.Log.Warning(str)
		s.r, s.err = strings.NewReader(str), err
		return s
	}

	s.Log.Info(fmt.Sprintf("Downloading %s source %s", s.area(), s.name))

	req.Header.Set("User-Agent", agent)
	if resp, err = (&http.Client{}).Do(req); err != nil {
		str := fmt.Sprintf("Unable to get response for %s", s.url)
		s.Log.Warning(str)
		s.r, s.err = strings.NewReader(str), err
		return s
	}

	body, err = ioutil.ReadAll(resp.Body)

	if len(body) < 1 {
		str := fmt.Sprintf("No data returned for %s", s.url)
		s.Log.Warning(str)
		s.r, s.err = strings.NewReader(str), err
		resp.Body.Close()
		return s
	}

	s.r, s.err = bytes.NewBuffer(body), err
	resp.Body.Close()
	return s
}
