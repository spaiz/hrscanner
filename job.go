package main

import "fmt"

// Job represents domain to be scanned, and returned with some filed data
type Job struct {
	ID     int
	Domain string
	Result string
	Error  error
}

func (r Job) String() string {
	return fmt.Sprintf("id: %d, domain: %s, header: %s", r.ID, r.Domain, r.Result)
}
