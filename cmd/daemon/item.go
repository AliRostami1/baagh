package main

import (
	"net/http"
)

type ItemGet struct {
}

func itemGetHandler(rw http.ResponseWriter, r *http.Request) {
	// core.GetItem()
}
func itemPostHandler(rw http.ResponseWriter, r *http.Request)   {}
func itemDeleteHandler(rw http.ResponseWriter, r *http.Request) {}
func itemWatchHandler(rw http.ResponseWriter, r *http.Request)  {}

func itemsGetHandler(rw http.ResponseWriter, r *http.Request)    {}
func itemsPostHandler(rw http.ResponseWriter, r *http.Request)   {}
func itemsDeleteHandler(rw http.ResponseWriter, r *http.Request) {}
func itemsWatchHandler(rw http.ResponseWriter, r *http.Request)  {}
