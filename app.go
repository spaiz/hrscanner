package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// NewApp injects dependencies and returns App instance to run
func NewApp(conf *Config) *App {
	resolver := NewDNSResolver()
	err := resolver.Load(conf.DNSServersFile)
	if err != nil {
		panic(err)
	}

	app := &App{
		Conf:  conf,
		Found: make([]*Job, 0),
		Client: &Client{
			DNSResolver: resolver,
		},
	}

	return app
}

// App runs fomain scrapping
type App struct {
	mu             sync.Mutex
	Client         *Client
	Conf           *Config
	Found          []*Job
	OnFound        func(job *Job)
	CompletedCount int
	FailedCount    int
	startTime      time.Time
}

// Run starts the workers
func (r *App) Run() error {
	logger.
		WithField("workers", r.Conf.WorkersNum).
		Debug("app started")

	var workersWait sync.WaitGroup
	var resultsWait sync.WaitGroup

	r.startTime = time.Now()
	jobsBuff := make(chan *Job, r.Conf.BufferSize)

	go LoadDomains(r.Conf.DomainsFile, jobsBuff)
	queueCompleted := make(chan *Job, r.Conf.BufferSize)

	resultsWait.Add(1)
	go func() {
		for job := range queueCompleted {
			r.CompletedCount++

			if job.Error != nil {
				logger.
					WithError(job.Error).
					WithField("domain", job.Domain).
					Error("Failed to scan the domain")

				r.FailedCount++
			}

			if job.Result != "" {
				r.Found = append(r.Found, job)
				if r.OnFound != nil {
					r.OnFound(job)
				}
			}
		}

		resultsWait.Done()
	}()

	workersWait.Add(r.Conf.WorkersNum)
	for i := 0; i < r.Conf.WorkersNum; i++ {
		worker := &Worker{
			JobsQueue:   jobsBuff,
			ID:          i,
			Client:      r.Client,
			Completed:   queueCompleted,
			WorkersWait: &workersWait,
		}

		go worker.Start()
	}

	workersWait.Wait()
	close(queueCompleted)
	resultsWait.Wait()

	return nil
}

// Freq returns scanning frequency
func (r *App) Freq() string {
	now := time.Now()
	duration := now.Sub(r.startTime)
	sec := duration.Seconds()
	if sec == 0 {
		sec = 1.0
	}

	freq := float64(r.CompletedCount) / sec
	return fmt.Sprintf("%.2f req/s", freq)
}

// FoundCount returns number of found headers
func (r *App) FoundCount() int {
	return len(r.Found)
}

// LoadDomains provides jobs by reading domains from the file.
// It will block when the buffer is full.
// It helps as to use a small memory amount.
func LoadDomains(file string, jobsBuff chan *Job) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	sc := bufio.NewScanner(f)
	i := 0
	for sc.Scan() {
		job := &Job{
			ID:     i,
			Domain: sc.Text(),
		}
		jobsBuff <- job
		i++
	}

	err = sc.Err()
	if err != nil  {
		logger.
			WithError(err).
			WithField("file", file).
			Error("there was some error on domains file scan")
	}

	close(jobsBuff)
}