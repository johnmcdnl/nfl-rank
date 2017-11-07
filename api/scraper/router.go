package scraper

import (
	"github.com/go-chi/chi"
	"net/http"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
	"strings"
	"github.com/sirupsen/logrus"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handler)
	return r
}

const (
	seasonParam = "season"
	phaseParam  = "phase"
	weekParam   = "week"
)

func handler(w http.ResponseWriter, r *http.Request) {

	season, err := getParam(r, seasonParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	phase, err := getParam(r, phaseParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	week, err := getParam(r, weekParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)

	go func(){
		if err := download( season, strings.ToUpper(phase), week); err != nil {
			logrus.Error(err)
			return
		}
		if err := writeScoreStrip(season, strings.ToUpper(phase), week); err != nil {
			logrus.Error(err)
			return
		}
	}()

}

func readOnlyDB() (*scrapeDb, error) {
	db, err := bolt.Open(nflDatabase, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, err
	}

	var dbCOnn scrapeDb
	dbCOnn.db = db
	return &dbCOnn, nil
}

func writableDB() (*scrapeDb, error) {
	db, err := bolt.Open(nflDatabase, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, err
	}

	var dbCOnn scrapeDb
	dbCOnn.db = db
	return &dbCOnn, nil
}

func getParam(r *http.Request, name string) (string, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return "", fmt.Errorf("param %s not found", name)
	}
	return value, nil
}
