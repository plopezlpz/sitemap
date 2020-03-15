package main

import (
	"crawler/internal/crawl"
	"flag"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	url := flag.String("url", "", "Base URL of the domain to be crawled")
	out := flag.String("out", "", "Path to the output file")
	logLevel := flag.String("logging", "info", "Sets the log level to debug")
	reqTimeoutSecs := flag.Int("timeout", 5, "Single HTTP request timeout in seconds")
	flag.Parse()

	setupLogging(logLevel)
	validateFlags(url, out)

	file := setupOutFile(out)
	defer file.Close()
	linksFetcher := setupFetcher(url, reqTimeoutSecs)

	start := time.Now()

	log.Info().Msgf("Crawling %v, this might take minutes...", *url)
	siteMap := <-crawl.GenerateSiteMap(*url, linksFetcher)

	log.Info().Msgf("Generating site map...")
	siteMap.PrintToFile(file)

	log.Info().Msgf("Finished site map of %v in %s", *url, time.Since(start))
}
