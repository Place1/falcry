package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type FalcoWebhook struct {
	Sink *Broadcaster
}

type FalcoEventChannel = chan *FalcoEvent

func NewFalcoWebhook(b *Broadcaster) *FalcoWebhook {
	return &FalcoWebhook{b}
}

type FalcoEvent struct {
	Output       string                 `json:"output"`
	Priority     string                 `json:"priority"`
	Rule         string                 `json:"rule"`
	Time         time.Time              `json:"time"`
	OutputFields FalcoEventOutputFields `json:"output_fields"`
	Raw          string                 `json:"-"`
}

type FalcoEventOutputFields = map[string]interface{}

func (f *FalcoWebhook) ListenAndServe(address string) error {
	router := mux.NewRouter()
	router.HandleFunc("/", f.webhook).Methods("POST")

	handler := negroni.New()
	handler.Use(negroni.NewRecovery())
	handler.UseHandler(router)

	server := http.Server{
		Addr:    address,
		Handler: handler,
	}

	logrus.Infof("falco webhook listening on %s", address)
	return server.ListenAndServe()
}

func (f *FalcoWebhook) webhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(500)
		return
	}

	event := &FalcoEvent{}
	if err := json.Unmarshal(body, event); err != nil {
		logrus.Error(err)
		w.WriteHeader(500)
		return
	}

	event.Raw = string(body)

	f.Sink.Send(event)

	w.WriteHeader(200)
}
