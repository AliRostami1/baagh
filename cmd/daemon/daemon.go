package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/internal/server"
	"github.com/AliRostami1/baagh/pkg/logy"
	flag "github.com/spf13/pflag"
)

func main() {
	port := flag.IntP("port", "p", 8080, "port of the web server")
	addr := flag.StringP("address", "a", "", "address of the web server")
	wto := flag.Duration("writeTimeout", 15*time.Second, "write timeout in seconds")
	rto := flag.Duration("readTimeout", 15*time.Second, "read timeout in seconds")
	ito := flag.Duration("idleTimeout", 60*time.Second, "idle timeout in seconds")
	shutdownWait := flag.Duration("shutdownWait", 10*time.Second, "time to shutdown after SIGINT")

	var logLevel logy.Level = logy.InfoLevel
	flag.VarP(&logLevel, "logLevel", "l", "set the log level")

	flag.Parse()

	app, err := application.New(logLevel)
	if err != nil {
		log.Fatalf("error happend during application startup: %v", err)
	}

	s := server.New(app.Ctx, app.Log, fmt.Sprintf("%s:%d", *addr, *port), *wto, *rto, *ito, *shutdownWait)

	err = s.Start()
	if err != nil {
		app.Log.Errorf("problem shutting down: %v", err)
		os.Exit(1)
	} else {
		app.Log.Infof("shutting down")
		os.Exit(0)
	}
}
