package main

import "sync"

// Worker will actually do domain scanning
type Worker struct {
	ID          int
	JobsQueue   chan *Job
	Completed   chan *Job
	Client      *Client
	WorkersWait *sync.WaitGroup
}

// Start starts worker to listen for jobs
func (r *Worker) Start() {
	defer r.WorkersWait.Done()

	for job := range r.JobsQueue {
		headers, err := r.Client.GetHeaders(job.Domain)
		job.Error = err
		if err == nil {
			job.Result = headers.Get("X-Recruiting")
		}

		r.Completed <- job
	}
}
