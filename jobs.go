// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	c "github.com/britannic/blacklist/config"
)

// getBlacklists assembles the http download jobs
func getBlacklists(timeout time.Duration, e excludes, urls []*c.Src) {
	jobs := make(chan Job, cores)
	results := make(chan Result, len(urls))
	done := make(chan struct{}, cores)

	go addJobs(jobs, urls, results)
	for i := 0; i < cores; i++ {
		go doJobs(done, jobs)
	}
	processResults(timeout, e, done, results)
}

// Result holds returned data
type Result struct {
	Data  string
	Error error
	Src   *c.Src
}

// processResults runs the data mill on the http content and writes it to
// its corresponding file
func processResults(timeout time.Duration, e excludes, done <-chan struct{},
	results <-chan Result) {
	finish := time.After(time.Duration(timeout))

	for working := cores; working > 0; {
		select {
		case result := <-results:
			data := process(result.Src, e, string(result.Data))
			fn := fmt.Sprintf("%v/%v.%v.blacklist.conf", dmsqDir, result.Src.Type, result.Src.Name)
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
			data := process(result.Src, e, string(result.Data))
			fn := fmt.Sprintf("%v/%v.%v.blacklist.conf", dmsqDir, result.Src.Type, result.Src.Name)
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

// Job holds job information
type Job struct {
	results chan<- Result
	src     *c.Src
}

// do is the pre-configured and http content loader
func (job Job) do() {
	if job.src.Name == "pre-configured" {
		var body string
		for key := range job.src.List {
			body += fmt.Sprintf("%v\n", key)
		}
		job.results <- Result{Src: job.src, Data: body, Error: nil}
	} else {
		const agent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
		client := new(http.Client)

		req, err := http.NewRequest("GET", job.src.URL, nil)
		if err != nil {
			log.Printf("Unable to form request for %s, error: %v", job.src.URL, err)
		}

		req.Header.Set("User-Agent", agent)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Unable to download %s, error: %v", job.src.URL, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
		}
		// log.Printf("Got results for job[%v]: (%v) %v\n", job.src.No, job.src.Type, job.src.Name)
		job.results <- Result{Src: job.src, Data: string(body[:]), Error: err}
	}
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
