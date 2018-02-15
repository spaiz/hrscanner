package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

var (
	workersNum     int
	bufferSize     int
	domainsFile    string
	dnsServersFile string
	resultsFile    string
	logsLevel      string

	logger = log.New()
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.IntVar(&workersNum, "workers", 250, "Number of workers to be created.")
	flag.IntVar(&bufferSize, "buffer", 1000, "Max domains number to be read from the file. Less number will lead to less memory usage.")
	flag.StringVar(&domainsFile, "domains-file", "./data/uniq-domains.txt", "Path to the file with domains separated by new line.")
	flag.StringVar(&dnsServersFile, "dns-file", "./data/dns-servers.txt", "Path to the file with dns servers IPs separated by new line.")
	flag.StringVar(&resultsFile, "results-file", "./results.txt", "Path to the file where found headers will be saved.")
	flag.StringVar(&logsLevel, "logs-level", "debug", "Logs level. (error, debug, ...)")
	flag.Parse()

	logger.Formatter = &log.JSONFormatter{}

	logger.Out = os.Stdout
	level, _ := log.ParseLevel(logsLevel)
	logger.SetLevel(level)
}

func main() {
	conf := &Config{
		WorkersNum:     workersNum,
		BufferSize:     bufferSize,
		DomainsFile:    domainsFile,
		DNSServersFile: dnsServersFile,
	}

	// open file for writing found headers
	file, err := os.OpenFile(resultsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	app := NewApp(conf)
	app.OnFound = func(job *Job) {
		logger.WithField("job", job).Debug("header found")
		file.WriteString(fmt.Sprintln(job))
	}

	// print summary every few second
	go func() {
		for {
			<-time.After(5 * time.Second)
			report(app)
		}
	}()

	err = app.Run()
	if err != nil {
		logger.WithError(err).Error("app finished")
	}

	report(app)
}

// MB used for user friendly memory usage printing
const MB uint64 = 1024 * 1024

// summary report to be written to the logs
func report(app *App) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory := struct {
		Alloc      string `json:"alloc"`
		TotalAlloc string `json:"total_alloc"`
		Sys        string `json:"sys"`
		NumGC      uint32 `json:"num_gc"`
	}{
		fmt.Sprintf("%d MB", m.Alloc/MB),
		fmt.Sprintf("%d MB", m.TotalAlloc/MB),
		fmt.Sprintf("%d MB", m.Sys/MB),
		m.NumGC,
	}

	logger.
		WithField("completed", app.CompletedCount).
		WithField("failed", app.FailedCount).
		WithField("found", app.FoundCount()).
		WithField("domains_found", app.Found).
		WithField("frequency", app.Freq()).
		WithField("memory", memory).
		Debug("report")
}
