package common

import (
	"net/http"
	"strings"
)

type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type WebHandle func(rw http.ResponseWriter, req *http.Request)

type Filter struct {
	// store the URIs that should be intercepted
	filterMap map[string]FilterHandle
}

func NewFiler() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandler(uri string) FilterHandle {
	return f.filterMap[uri]
}

func (f *Filter) Handle(handle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) {
				// intercepted the uri and activate the handle
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}

		handle(rw, r)
	}
}
