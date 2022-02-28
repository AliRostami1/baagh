package server

import (
	"fmt"
	"net/http"

	"github.com/AliRostami1/baagh/internal/version"
)

type Version struct {
	ApplicationVersion string `json:"applicationVersion"`
	GpiodVersion       string `json:"gpiodVersion"`
	UAPIVersion        string `json:"uApiVersion"`
	KernelVersion      string `json:"kernelVersion"`
}

func versionHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, version.Version())
}
