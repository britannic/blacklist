// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/Sirupsen/logrus"

	c "github.com/britannic/blacklist/config"
)

// getBlacklists assembles the http download jobs
func getBlacklists(timeout time.Duration, d c.Dict, e c.Dict, a areaURLs) {

	for k := range a {
		jobs := make(chan Job, cores)
		results := make(chan Result, len(a[k]))
		done := make(chan struct{}, cores)

		go addJobs(jobs, a[k], results)
		for i := 0; i < cores; i++ {
			go doJobs(done, jobs)
		}
		processResults(timeout, d, e, done, results)
	}
}

// Result holds returned data
type Result struct {
	Data  string
	Error error
	Src   *c.Src
}

// processResults mills the http content and writes it to its corresponding file
func processResults(timeout time.Duration, d c.Dict, e c.Dict, done <-chan struct{},
	results <-chan Result) {
	finish := time.After(time.Duration(timeout))

	for working := cores; working > 0; {
		select {
		case result := <-results:
			data := process(result.Src, d, e, string(result.Data))
			fn := fmt.Sprintf(fStr, dmsqDir, result.Src.Type, result.Src.Name)
			log.Printf("[Select 1] writing job[%v] %v\n", result.Src.No, fn)
			if err := writeFile(fn, getList(data)); err != nil {
				fmt.Println(err)
			}
		case <-finish:
			log.Println("timed out")
			return
		case <-done:
			working--
		}
	}
	for {
		select {
		case result := <-results:
			data := process(result.Src, d, e, string(result.Data))
			fn := fmt.Sprintf(fStr, dmsqDir, result.Src.Type, result.Src.Name)
			log.Printf("[Select 2] writing job[%v] %v\n", result.Src.No, fn)
			if err := writeFile(fn, getList(data)); err != nil {
				log.Println("Error: ", err)
			}
		case <-finish:
			log.Println("timed out")
			return
		default:
			return
		}
	}
}

func getHTTP(URL string) (body []byte, err error) {
	const agent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
	var (
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest("GET", URL, nil)
	if err == nil {
		req.Header.Set("User-Agent", agent)
		// req.Header.Add("Content-Type", "application/json")
		debug(httputil.DumpRequestOut(req, true))
		resp, err = (&http.Client{}).Do(req)
	} else {
		log.Printf("Unable to form request for %s, error: %v", URL, err)
	}

	if err == nil {
		defer resp.Body.Close()
		debug(httputil.DumpResponse(resp, true))
		body, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// Job holds job information
type Job struct {
	results chan<- Result
	src     *c.Src
}

// do is the pre-configured and http content loader
func (job Job) do() {
	var body []byte
	var err error
	switch {
	case job.src.Name == "pre-configured":
		for key := range job.src.List {
			body = append(body, fmt.Sprintf("%v\n", key)...)
		}
	default:
		body, err = getHTTP(job.src.URL)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}
	job.results <- Result{Src: job.src, Data: string(body[:]), Error: err}
}

// addJobs puts jobs on the tasklist
func addJobs(jobs chan<- Job, urls []*c.Src, results chan<- Result) {
	for i, url := range urls {
		url.No = i + 1
		jobs <- Job{src: url, results: results}
		// log.Printf("Adding download: (%v) %v\n", url.Type, url.Name)
	}
	close(jobs)
}

// doJobs supervises the jobs and uses an empty struct to signal completion
func doJobs(done chan<- struct{}, jobs <-chan Job) {
	for job := range jobs {
		job.do()
		// log.Printf("Running job[%v]: (%v) %v\n", job.src.No, job.src.Type, job.src.Name)
	}
	done <- struct{}{}
}

func debug(data []byte, err error) {
	if dbg == false {
		return
	}
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
