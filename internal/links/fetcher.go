package links

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type linkFetcher struct {
	client http.Client
	parser parser
}

type Fetcher interface {
	FetchLinks(url string) *PageLinks
}

func NewFetcher(site string, client http.Client) (Fetcher, error) {
	p, err := newParser(site)
	if err != nil {
		return &linkFetcher{}, err
	}
	return &linkFetcher{
		client: client,
		parser: p,
	}, nil
}

func (f *linkFetcher) FetchLinks(url string) *PageLinks {
	url = f.parser.toFullURL(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Err(err).Msgf("Error getting %v, skipping", url)
		return nil
	}
	log.Trace().Msgf("Parsing %v", url)
	defer resp.Body.Close()

	res := f.parser.parse(url, resp.Body)
	return &res
}
