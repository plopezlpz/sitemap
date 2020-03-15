package store

import (
	"reflect"
	"testing"
)

func TestSite_Print(t *testing.T) {
	tests := []struct {
		name  string
		store map[string]map[string]bool
	}{
		{
			name: "test printing",
			store: map[string]map[string]bool{
				"/":   {"a1": true, "a2": true, "a3": true},
				"a1":  {"/": true, "a2": true, "a11": true, "a12": true},
				"a11": {"a2": true, "a12": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				store: tt.store,
			}
			s.PrintToStdOut()
		})
	}
}

func TestSite_SaveLinks(t *testing.T) {
	type args struct {
		url   string
		links []string
	}
	tests := []struct {
		name   string
		store  map[string]map[string]bool
		args   args
		cbWant map[string]bool
		want   map[string]map[string]bool
	}{
		{
			name:  "add new links to store and callback",
			store: map[string]map[string]bool{},
			args: args{
				url:   "/a",
				links: []string{"/b", "/c"},
			},
			cbWant: map[string]bool{
				"/b": true, "/c": true,
			},
			want: map[string]map[string]bool{
				"/a": {"/b": true, "/c": true},
				"/b": {},
				"/c": {},
			},
		},
		{
			name: "skip adding existing links and don't callback",
			store: map[string]map[string]bool{
				"/b": {},
			},
			args: args{
				url:   "/a",
				links: []string{"/b", "/c"},
			},
			cbWant: map[string]bool{
				"/c": true,
			},
			want: map[string]map[string]bool{
				"/a": {"/b": true, "/c": true},
				"/b": {},
				"/c": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				store: tt.store,
			}
			cbActual := map[string]bool{}
			s.SaveLinks(tt.args.url, tt.args.links, func(url string) {
				cbActual[url] = true
			})

			if !reflect.DeepEqual(tt.store, tt.want) {
				t.Errorf("store = %v, want %v", tt.store, tt.want)
			}
			if !reflect.DeepEqual(cbActual, tt.cbWant) {
				t.Errorf("new link callback called with = %v, wanted = %v", cbActual, tt.cbWant)
			}
		})
	}
}
