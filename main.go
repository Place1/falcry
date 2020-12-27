package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/place1/falcry/internal/api"
	"github.com/place1/falcry/internal/storage"
	"github.com/place1/falcry/internal/web"
	"github.com/place1/falcry/internal/webhook"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("falcry", "A tiny webui for Falco")
	logLevel = app.Flag("log-level", "Log level: trace, debug, info, error, fatal").Default("info").String()
)

func main() {
	// parse CLI arguments
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// set global log level
	if level, err := logrus.ParseLevel(*logLevel); err == nil {
		logrus.SetLevel(level)
	}

	wg := sync.WaitGroup{}
	db := storage.MustOpen("sqlite3://db.sqlite3")
	broadcaster := webhook.NewBroadcaster()

	// Webhook listener
	wg.Add(1)
	go func() {
		webhook.NewFalcoWebhook(broadcaster).ListenAndServe("0.0.0.0:8000")
		wg.Done()
	}()

	// Storage event sink
	wg.Add(1)
	go func() {
		storage.NewStorageSink(db, broadcaster).ListenAndSave()
		wg.Done()
	}()

	// Website + API server
	wg.Add(1)
	go func() {
		router := mux.NewRouter()

		router.PathPrefix("/api").Handler(api.New(db, broadcaster).Handler())
		router.PathPrefix("/").Handler(web.New().Handler())

		server := http.Server{
			Addr:    "0.0.0.0:8000",
			Handler: router,
		}

		server.ListenAndServe()
		wg.Done()
	}()

	wg.Wait()
}
