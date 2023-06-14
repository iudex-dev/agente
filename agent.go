package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Agent struct {
	client http.Client
}

func (A *Agent) getPendingJobs() []string {
	res, err := A.client.Get(fmt.Sprintf("%s/api/jobs/pending", iudexApi))
	if err != nil {
		exitWithError(fmt.Sprintf("Error connecting to '%s': %s", iudexApi, err))
	}
	var result []string
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		exitWithError(fmt.Sprintf("Cannot decode result: %s", err))
	}
	res.Body.Close()
	return result
}

type Job struct {
	Id string `json:"id"`
	Code string `json:"code"`
}

func (A *Agent) getJob(id string) Job {
	res, err := A.client.Get(fmt.Sprintf("%s/api/jobs/%s", iudexApi, id))
	if err != nil {
		exitWithError(fmt.Sprintf("Error connecting to '%s': %s", iudexApi, err))
	}
	job := Job{}
	err = json.NewDecoder(res.Body).Decode(&job)
	if err != nil {
		exitWithError(fmt.Sprintf("Cannot decode result: %s", err))
	}
	res.Body.Close()
	return job
}

func (A *Agent) Run() {
	for {
		pending := A.getPendingJobs()
		job := A.getJob(pending[0]) // FIXME: Take one at random?
		result := CarcerResult{}
		err := ProcessSubmission(job.Code, &result)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
		} else {
			fmt.Println("Result", result)
		}

		time.Sleep(1000 * time.Millisecond) // FIXME: Wait at random?
	}
}