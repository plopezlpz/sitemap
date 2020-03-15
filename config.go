package main

import (
	"crawler/internal/links"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func validateFlags(url, out *string) {
	if *url == "" || *out == "" {
		log.Error().Msgf("missing required arguments")
		flag.Usage()
		os.Exit(1)
	}
}

func setupOutFile(out *string) *os.File {
	file, err := os.OpenFile(*out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not open file to print")
	}
	return file
}

func setupLogging(level *string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	l, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Msgf("invalid log level %v", *level)
	}
	zerolog.SetGlobalLevel(l)
}

func setupFetcher(url *string, reqTimeoutSecs *int) links.Fetcher {
	httpClient := http.Client{
		Timeout: time.Duration(*reqTimeoutSecs) * time.Second,
	}
	c, err := links.NewFetcher(*url, httpClient)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot setup client")
	}
	return c
}
