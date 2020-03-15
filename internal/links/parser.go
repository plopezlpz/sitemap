package links

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"io"
	"net/url"
)

type linkFilter func(link *string) *string

type PageLinks struct {
	Url   string
	Links []string
}

type parser struct {
	base url.URL
}

func newParser(site string) (parser, error) {
	l, err := url.Parse(site)
	if err != nil {
		return parser{}, err
	}
	return parser{base: *l}, nil
}

func (p *parser) toFullURL(uri string) string {
	l, err := p.base.Parse(uri)
	if err != nil {
		log.Err(err).Msgf("error parsing %v", uri)
		return uri
	}
	return l.String()
}

// parse the internal Links belonging to the same domain as base
func (p *parser) parse(pageUrl string, reader io.Reader) PageLinks {
	links := getLinks(reader, sameDomainFilter(p.base))
	return PageLinks{Url: requestURI(pageUrl), Links: links}
}

func sameDomainFilter(base url.URL) linkFilter {
	return func(link *string) *string {
		if link == nil {
			return nil
		}
		l, err := base.Parse(*link)
		if err != nil {
			log.Err(err).Msgf("skipping %v", *link)
			return nil
		}
		if base.Hostname() == l.Hostname() {
			uri := l.RequestURI()
			return &uri
		}
		return nil
	}
}

func getLinks(reader io.Reader, filter linkFilter) []string {
	var uniqueLinks = make(map[string]bool)
	tokenizer := html.NewTokenizer(reader)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}
		} else if tokenType == html.StartTagToken {
			href := filter(getHref(tokenizer.Token()))
			if href != nil {
				uniqueLinks[*href] = true
			}
		}
	}
	return toArr(uniqueLinks)
}

func toArr(keyMap map[string]bool) []string {
	keys := make([]string, len(keyMap))
	i := 0
	for k := range keyMap {
		keys[i] = k
		i++
	}
	return keys
}

func getHref(token html.Token) *string {
	if "a" == token.Data {
		// get the href attribute
		for _, att := range token.Attr {
			if att.Key == "href" {
				return &att.Val
			}
		}
	}
	return nil
}

func requestURI(link string) string {
	l, err := url.Parse(link)
	if err != nil {
		log.Err(err).Msgf("error parsing %v, skipping", link)
		return ""
	}
	return l.RequestURI()
}
