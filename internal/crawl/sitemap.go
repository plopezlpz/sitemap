package crawl

import (
	"crawler/internal/links"
	"crawler/internal/store"
	"github.com/rs/zerolog/log"
)

func GenerateSiteMap(url string, fetcher links.Fetcher) chan store.Site {
	siteMapCh := make(chan store.Site)
	go func() {
		siteMap := store.NewSite()
		linksCh := make(chan *links.PageLinks)
		count := 1
		go func() {
			linksCh <- fetcher.FetchLinks(url)
		}()

		for l := range linksCh {
			count--
			if l != nil {
				siteMap.SaveLinks(l.Url, l.Links, func(newLink string) {
					count++
					go func() {
						linksCh <- fetcher.FetchLinks(newLink)
					}()
				})
			}
			log.Trace().Msgf("Go routines %v", count)
			if count <= 0 {
				siteMapCh <- siteMap
			}
		}
	}()
	return siteMapCh
}
