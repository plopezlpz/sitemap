package crawl

import (
	"crawler/internal/links"
	"crawler/internal/mock"
	"reflect"
	"testing"
)

func TestGenerateSiteMap(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		want      map[string]map[string]bool
		mockFetch map[string]links.PageLinks
	}{
		{
			name: "generates a site map",
			url:  "/",
			mockFetch: map[string]links.PageLinks{
				"/": {
					Url:   "/",
					Links: []string{"/a", "/b"},
				},
				"/b": {
					Url:   "/b",
					Links: []string{"/a", "/c"},
				},
				"/a": {
					Url:   "/a",
					Links: []string{"/b", "/"},
				},
				"/c": {
					Url:   "/c",
					Links: []string{},
				},
			},
			want: map[string]map[string]bool{
				"/": {
					"/a": true, "/b": true,
				},
				"/a": {
					"/b": true, "/": true,
				},
				"/b": {
					"/a": true, "/c": true,
				},
				"/c": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := mock.NewFetcher(tt.mockFetch)
			if got := <-GenerateSiteMap(tt.url, fetcher); !reflect.DeepEqual(got.Get(), tt.want) {
				t.Errorf("GenerateSiteMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
