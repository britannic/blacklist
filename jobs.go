package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	c "github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
)

// getBlacklists assembles the http download jobs
func getBlacklists(timeout time.Duration, dex c.Dict, ex c.Dict, a data.AreaURLs) {

	for k := range a {
		jobs := make(chan Job, cores)
		results := make(chan Result, len(a[k]))
		done := make(chan struct{}, cores)

		go addJobs(jobs, a[k], results)
		for i := 0; i < cores; i++ {
			go doJobs(done, jobs)
		}
		getResults(timeout, dex, ex, done, results)
	}
}

// Result holds returned data
type Result struct {
	Data  []byte
	Error error
	Src   *c.Src
}

// getResults collects the HTTP content and sends it to processResults
func getResults(timeout time.Duration, dex c.Dict, ex c.Dict, done <-chan struct{},
	results <-chan Result) {
	finish := time.After(time.Duration(timeout))

	for working := cores; working > 0; {
		select {
		case result := <-results:
			if err := processResults(&result, dex, ex); err != nil {
				log.Println("Error: ", err)
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
			if err := processResults(&result, dex, ex); err != nil {
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

// processResults mills the http content and writes it to its corresponding file
func processResults(result *Result, dex c.Dict, ex c.Dict) (err error) {
	b := bufio.NewScanner(strings.NewReader(string(result.Data)))
	pdata := data.Process(result.Src, dex, ex, b)
	fn := fmt.Sprintf(g.FStr, g.DmsqDir, result.Src.Type, result.Src.Name)
	log.Printf("Writing job[%v] %v\n", result.Src.No, fn)
	// if result.Src.Name == "pre-configured" {
	// fmt.Println(string(result.Data))
	// 	fmt.Println(pdata)
	// }
	err = utils.WriteFile(fn, data.GetList(pdata))
	return err
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
		body, err = data.GetHTTP(job.src.URL)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}
	job.results <- Result{Src: job.src, Data: body, Error: err}
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
