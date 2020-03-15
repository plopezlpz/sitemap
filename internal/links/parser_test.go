package links

import (
	"io"
	"log"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func toUrl(urlStr string) url.URL {
	base, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("error parsing Url in test %v", urlStr)
	}
	return *base
}

func TestLinkParser_ExtractLinks(t *testing.T) {
	type args struct {
		reader io.Reader
		url    string
	}
	tests := []struct {
		name    string
		baseUrl url.URL
		args    args
		want    PageLinks
	}{
		{
			name:    "should extract Links belonging to the base domain",
			baseUrl: toUrl("http://mine.com"),
			args: args{reader: strings.NewReader(`
				<html>
					<a href="http://mine.com">home</a>
					<a href="/bills">mine</a>
					<a href="/bills/one">mine</a>
				</html>
			`)},
			want: PageLinks{
				Url:   "/",
				Links: []string{"/", "/bills", "/bills/one"},
			},
		},
		{
			name:    "should skip Links not belonging to the base domain",
			baseUrl: toUrl("http://mine.com"),
			args: args{reader: strings.NewReader(`
				<html>
					<a href="http://mine.com/home">home</a>
					<a href="http://other.mine.com/other">home</a>
					<a href="/bills">mine</a>
					<a href="/bills/one">mine</a>
				</html>
			`)},
			want: PageLinks{
				Url:   "/",
				Links: []string{"/home", "/bills", "/bills/one"},
			},
		},
		{
			name:    "should include query string",
			baseUrl: toUrl("http://mine.com"),
			args: args{reader: strings.NewReader(`
				<html>
					<a href="/bills">mine</a>
					<a href="/bills/one?hello=true">mine</a>
				</html>
			`)},
			want: PageLinks{
				Url:   "/",
				Links: []string{"/bills", "/bills/one?hello=true"},
			},
		},
		{
			name:    "should not include fragment string",
			baseUrl: toUrl("http://mine.com"),
			args: args{reader: strings.NewReader(`
				<html>
					<a href="/bills">mine</a>
					<a href="/bills/one#hello=true">mine</a>
				</html>
			`)},
			want: PageLinks{
				Url:   "/",
				Links: []string{"/bills", "/bills/one"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				base: tt.baseUrl,
			}
			got := p.parse(tt.args.url, tt.args.reader)
			sort.Strings(got.Links)
			sort.Strings(tt.want.Links)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_toFullURL(t *testing.T) {
	type fields struct {
		base url.URL
	}
	type args struct {
		uri string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should return full url with starting with /",
			fields: fields{
				base: toUrl("http://mine.com"),
			},
			args: args{
				uri: "/a",
			},
			want: "http://mine.com/a",
		},
		{
			name: "should return full url with starting without /",
			fields: fields{
				base: toUrl("http://mine.com"),
			},
			args: args{
				uri: "a",
			},
			want: "http://mine.com/a",
		},
		{
			name: "should return full url with query params",
			fields: fields{
				base: toUrl("http://mine.com"),
			},
			args: args{
				uri: "/a?a=b",
			},
			want: "http://mine.com/a?a=b",
		},
		{
			name: "should return full url with query params",
			fields: fields{
				base: toUrl("http://mine.com"),
			},
			args: args{
				uri: "http://mine.com/a",
			},
			want: "http://mine.com/a",
		},
		{
			name: "should return full url with query params",
			fields: fields{
				base: toUrl("http://mine.com"),
			},
			args: args{
				uri: "http://yours.com/a",
			},
			want: "http://yours.com/a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				base: tt.fields.base,
			}
			if got := p.toFullURL(tt.args.uri); got != tt.want {
				t.Errorf("toFullURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
