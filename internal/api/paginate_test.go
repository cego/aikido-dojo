package api

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

type page struct {
	body    string
	hasNext string // X-Has-Next-Page value; "" sends no header
	status  int
}

// pageServer serves pages[N] for ?page=N and records every query string.
type pageServer struct {
	mu      sync.Mutex
	queries []string
	pages   []page
}

func (s *pageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.queries = append(s.queries, r.URL.RawQuery)
	s.mu.Unlock()
	n, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if n >= len(s.pages) {
		http.Error(w, "no such page", http.StatusNotFound)
		return
	}
	p := s.pages[n]
	if p.hasNext != "" {
		w.Header().Set("X-Has-Next-Page", p.hasNext)
	}
	if p.status != 0 {
		w.WriteHeader(p.status)
	}
	fmt.Fprint(w, p.body)
}

func (s *pageServer) seen() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.queries)
}

func collect(seq iter.Seq2[json.RawMessage, error]) ([]string, error) {
	var got []string
	for item, err := range seq {
		if err != nil {
			return got, err
		}
		got = append(got, string(item))
	}
	return got, nil
}

func TestItems(t *testing.T) {
	arrays := Paging{Style: PageArray, SizeParam: "per_page", Size: 2}
	tests := []struct {
		name       string
		paging     Paging
		pages      []page
		want       []string
		wantCalls  int
		wantFirstQ string
	}{
		{
			name: "an array list ends on an empty page", paging: arrays,
			pages: []page{{body: `[1,2]`}, {body: `[3]`}, {body: `[]`}},
			want:  []string{"1", "2", "3"}, wantCalls: 3, wantFirstQ: "filter_x=a&page=0&per_page=2",
		},
		{
			name: "a header list keeps going past an empty page", paging: Paging{Style: PageHeader, SizeParam: "per_page", Size: 100},
			pages: []page{{body: `[1]`, hasNext: "true"}, {body: `[]`, hasNext: "true"}, {body: `[2]`, hasNext: "false"}},
			want:  []string{"1", "2"}, wantCalls: 3, wantFirstQ: "filter_x=a&page=0&per_page=100",
		},
		{
			name: "a header list stops when the header is missing", paging: Paging{Style: PageHeader, SizeParam: "per_page", Size: 100},
			pages: []page{{body: `[1]`}}, want: []string{"1"}, wantCalls: 1, wantFirstQ: "filter_x=a&page=0&per_page=100",
		},
		{
			name:   "an envelope list follows hasMore",
			paging: Paging{Style: PageEnvelope, SizeParam: "limit", Size: 50, Items: "assets", More: "hasMore"},
			pages:  []page{{body: `{"assets":[1],"hasMore":true,"totalCount":2}`}, {body: `{"assets":[2],"hasMore":false,"totalCount":2}`}},
			want:   []string{"1", "2"}, wantCalls: 2, wantFirstQ: "filter_x=a&limit=50&page=0",
		},
		{
			name:   "an envelope list without hasMore ends on an empty page",
			paging: Paging{Style: PageEnvelope, SizeParam: "per_page", Size: 50, Items: "users"},
			pages:  []page{{body: `{"users":[1]}`}, {body: `{"users":[]}`}},
			want:   []string{"1"}, wantCalls: 2, wantFirstQ: "filter_x=a&page=0&per_page=50",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &pageServer{pages: tt.pages}
			ta := newTestAPI(t, srv.ServeHTTP, nil)
			req := Request{Method: http.MethodGet, Path: "/list", Query: map[string][]string{"filter_x": {"a"}}}
			got, err := collect(ta.client.Items(t.Context(), req, tt.paging))
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("items = %v, want %v", got, tt.want)
			}
			if q := srv.seen(); len(q) != tt.wantCalls || q[0] != tt.wantFirstQ {
				t.Errorf("queries = %q; want %d calls starting with %q", q, tt.wantCalls, tt.wantFirstQ)
			}
			if len(req.Query) != 1 {
				t.Errorf("Items changed the caller's query: %v", req.Query)
			}
		})
	}
}

func TestItemsStopsFetchingWhenTheCallerStops(t *testing.T) {
	srv := &pageServer{pages: []page{{body: `[1,2]`}, {body: `[3,4]`}}}
	ta := newTestAPI(t, srv.ServeHTTP, nil)
	for range ta.client.Items(t.Context(), Request{Method: http.MethodGet, Path: "/list"}, Paging{Style: PageArray, SizeParam: "per_page", Size: 2}) {
		break
	}
	if got := len(srv.seen()); got != 1 {
		t.Errorf("calls = %d, want 1: taking one item must not fetch the next page", got)
	}
}

func TestItemsFailures(t *testing.T) {
	tests := []struct {
		name     string
		paging   Paging
		pages    []page
		want     []string
		wantText string
	}{
		{
			name: "an API error after the first page keeps the items already yielded", paging: Paging{Style: PageArray, SizeParam: "per_page", Size: 1},
			pages: []page{{body: `[1]`}, {body: `{"error":"gone"}`, status: 404}}, want: []string{"1"}, wantText: "gone",
		},
		{
			name:   "an envelope without its items field",
			paging: Paging{Style: PageEnvelope, SizeParam: "limit", Size: 1, Items: "assets"},
			pages:  []page{{body: `{"things":[]}`}}, wantText: `no "assets" field`,
		},
		{
			name:   "a hasMore that is not a boolean",
			paging: Paging{Style: PageEnvelope, SizeParam: "limit", Size: 1, Items: "assets", More: "hasMore"},
			pages:  []page{{body: `{"assets":[1],"hasMore":"yes"}`}}, wantText: `"hasMore"`,
		},
		{
			name:   "an envelope body that is not an object",
			paging: Paging{Style: PageEnvelope, SizeParam: "limit", Size: 1, Items: "assets"},
			pages:  []page{{body: `[1]`}}, wantText: "want a JSON object",
		},
		{
			name:   "an envelope whose items are not an array",
			paging: Paging{Style: PageEnvelope, SizeParam: "limit", Size: 1, Items: "assets"},
			pages:  []page{{body: `{"assets":{}}`}}, wantText: `field "assets"`,
		},
		{
			name: "a body that is not an array", paging: Paging{Style: PageArray, SizeParam: "per_page", Size: 1},
			pages: []page{{body: `{"items":[]}`}}, wantText: "decode page 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &pageServer{pages: tt.pages}
			ta := newTestAPI(t, srv.ServeHTTP, nil)
			got, err := collect(ta.client.Items(t.Context(), Request{Method: http.MethodGet, Path: "/list"}, tt.paging))
			if err == nil || !strings.Contains(err.Error(), tt.wantText) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantText)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("items before the error = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("the API error keeps its exit code", func(t *testing.T) {
		srv := &pageServer{pages: []page{{body: `{}`, status: 404}}}
		ta := newTestAPI(t, srv.ServeHTTP, nil)
		_, err := collect(ta.client.Items(t.Context(), Request{Method: http.MethodGet, Path: "/list"}, Paging{Style: PageArray, SizeParam: "per_page", Size: 1}))
		wantCode(t, err, "not_found", clierr.ExitNotFound)
	})
}
