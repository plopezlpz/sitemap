package store

type Site struct {
	store map[string]map[string]bool
}

func NewSite() Site {
	return Site{store: make(map[string]map[string]bool)}
}

// SaveLinks and run the `onNewFn` on each
// new link that has not previously been saved
func (s *Site) SaveLinks(url string, links []string, onNewFn func(url string)) {

	if s.store[url] == nil {
		s.store[url] = make(map[string]bool)
	}

	for _, l := range links {
		s.store[url][l] = true
		if s.store[l] == nil {
			s.store[l] = make(map[string]bool)

			// call the onNewFn function for not previously seen links
			onNewFn(l)
		}
	}
}

func (s *Site) Get() map[string]map[string]bool {
	return s.store
}
