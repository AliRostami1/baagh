package main

import (
	"fmt"
	"net/http"
)

func versionHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, version())
}
