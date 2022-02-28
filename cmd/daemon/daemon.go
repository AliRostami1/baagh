package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/logy"
	"github.com/gorilla/mux"
	flag "github.com/spf13/pflag"
)

func version() string {
	return "0.0.1"
}

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

	log.Printf("log level: %s", logLevel.String())
	log.Printf("port: %d", *port)

	app, err := application.New(logLevel)
	if err != nil {
		log.Fatalf("error happend during application startup: %v", err)
	}

	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, version())
	})

	// apiSub := r.PathPrefix("/api").Subrouter()

	// chipSub := apiSub.PathPrefix("/chip").Subrouter()

	// chipSub.HandleFunc("/item", itemsGetHandler).Methods("GET")
	// chipSub.HandleFunc("/item", itemsWatchHandler).Methods("GET").Queries("watch", "true")
	// chipSub.HandleFunc("/item", itemsPostHandler).Methods("POST")
	// chipSub.HandleFunc("/item", itemsDeleteHandler).Methods("DELETE")

	// chipSub.HandleFunc("/item/{offset:[0-9]+}", itemGetHandler).Methods("GET")
	// chipSub.HandleFunc("/item/{offset:[0-9]+}", itemWatchHandler).Methods("GET").Queries("watch", "true")
	// chipSub.HandleFunc("/item/{offset:[0-9]+}", itemPostHandler).Methods("POST")
	// chipSub.HandleFunc("/item/{offset:[0-9]+}", itemDeleteHandler).Methods("DELETE")

	// apiSub.HandleFunc("/healthcheck", healthCheckHandler).Methods("GET")

	// apiSub.HandleFunc("/version", func(rw http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(rw, version())
	// }).Methods("GET")

	srv := &http.Server{
		// Addr: fmt.Sprintf("%s:%d", *addr, *port),
		Addr: ":8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * *wto,
		ReadTimeout:  time.Second * *rto,
		IdleTimeout:  time.Second * *ito,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		app.Log.Infof("starting the server on %s:%d", *addr, *port)
		if err := srv.ListenAndServe(); err == http.ErrServerClosed {
			app.Log.Infof("server stopped: %v", err)
		} else {
			app.Log.Errorf("server crashed: %v", err)
		}
	}()

	<-app.Ctx.Done()

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *shutdownWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	app.Log.Infof("shutting down")
	os.Exit(0)
}
