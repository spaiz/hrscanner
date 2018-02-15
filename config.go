package main

// Config used to be more flexible
type Config struct {
	WorkersNum     int
	BufferSize     int
	DomainsFile    string
	DNSServersFile string
}
