package mock

import "crawler/internal/links"

type mockFetcher struct {
	mockedCalls map[string]links.PageLinks
}

func NewFetcher(mockedCalls map[string]links.PageLinks) links.Fetcher {
	return &mockFetcher{
		mockedCalls: mockedCalls,
	}
}

func (f *mockFetcher) FetchLinks(url string) *links.PageLinks {
	res := f.mockedCalls[url]
	return &res
}
