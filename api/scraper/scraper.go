package scraper

import (
	"github.com/boltdb/bolt"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/sirupsen/logrus"
	"encoding/xml"
	"github.com/johnmcdnl/nfl-rank/nfl"
	"strings"
	"encoding/json"
)

const (
	nflDatabase         = "nfl-data.db"
	nflRawXML           = "nfl-raw-xml"
	nflScoreStripBucket = "nfl-scorestrip-json"
)

type scrapeDb struct {
	db *bolt.DB
}

func download(season, phase, week string) error {
	const scrapeURLFormat = "http://www.nfl.com/ajax/scorestrip?season=%s&seasonType=%s&week=%s"
	resp, err := http.Get(fmt.Sprintf(scrapeURLFormat, season, phase, week))
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}


	db, err := writableDB()
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer db.db.Close()

	if err := deleteRecord(db, season, phase, week); err != nil {
		return err
	}
	return writeRecord(db, season, phase, week, body)
}

func writeRecord(db *scrapeDb, season, phase, week string, data []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nflRawXML))
		if err != nil {
			logrus.Error(err)
			return err
		}
		logrus.Infoln("Creating ", string(name(season, phase, week)))
		return bucket.Put(name(season, phase, week), data)
	})
}

func deleteRecord(db *scrapeDb, season, phase, week string) error {
	logrus.Infoln("Deleting: ", string(name(season, phase, week)))
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nflRawXML))
		if err != nil {
			logrus.Error(err)
			return err
		}
		return bucket.Delete(name(season, phase, week))
	})
}

func parse(season, phase, week string) (*nfl.ScoreStrip, error) {

	db, err := readOnlyDB()
	if err != nil {
		return nil, err
	}
	defer db.db.Close()

	var s nfl.ScoreStrip
	if err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(nflRawXML))
		data := bucket.Get(name(season, phase, week))
		return xml.Unmarshal(data, &s.GameWeek)
	}); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &s, nil
}

func writeScoreStrip(season, phase, week string) error {
	scoreStrip, err := parse(season, strings.ToUpper(phase), week)
	if err != nil {
		return err
	}

	scoreStripJson, err := json.Marshal(scoreStrip)
	if err != nil {
		return err
	}

	db, err := writableDB()
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer db.db.Close()
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nflScoreStripBucket))
		if err != nil {
			logrus.Error(err)
			return err
		}

		logrus.Infoln("Creating JSON", string(name(season, phase, week)))
		return bucket.Put(name(season, phase, week), scoreStripJson)
	})
}

func name(season, phase, week string) []byte {
	return []byte(fmt.Sprintf("%s_%s_%s", season, phase, week))
}
