package nfl

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/boltdb/bolt"
	"time"
	"encoding/xml"
	"log"
)

const (
	firstSeason = 1900
	lastSeason  = 2017

	firstWeek = 0
	lastWeek  = 30

	currentWeek = 5

	PreSeason     = "PRE"
	RegularSeason = "REG"
	PostSeason    = "POST"

	nflDatabase = "nfl-data.db"
	nflRawXML   = "nfl-raw-xml"
)

type DB struct {
	DB *bolt.DB
}

func ScrapeAll() []*ScoreStrip {
	db, err := bolt.Open(nflDatabase, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var dbCOnn DB
	dbCOnn.DB = db
	setup(&dbCOnn)
	downloadAll(&dbCOnn)
	return parseAll(&dbCOnn)
}

func setup(db *DB) {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(nflRawXML))
		return err
	})
	if err != nil {
		panic(err)
	}
}

func parse(db *DB, season int, phase string, week int) *ScoreStrip {
	var s ScoreStrip
	if err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(nflRawXML))
		data := bucket.Get(name(season, phase, week))
		return xml.Unmarshal(data, &s.GameWeek)
	}); err != nil {
		deleteRecord(db, season, phase, week)
		panic(err)
	}
	return &s
}

func deleteRecord(db *DB, season int, phase string, week int) {
	log.Println("DELETE: ", string(name(season, phase, week)))
	if err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nflRawXML))
		if err != nil {
			return err
		}
		return bucket.Delete(name(season, phase, week))
	}); err != nil {
		panic(err)
	}
}

func parseAll(dbCOnn *DB) []*ScoreStrip {
	var scoreStrips []*ScoreStrip
	for season := firstSeason; season <= lastSeason; season++ {
		for week := firstWeek; week <= lastWeek; week++ {
			if season == 2017 {
				ss := parse(dbCOnn, season, PreSeason, week)
				if ss.GameWeek.Games != nil {
					scoreStrips = append(scoreStrips, ss)
				}

				if week <= currentWeek {
					ss := parse(dbCOnn, season, RegularSeason, week)
					if ss.GameWeek.Games != nil {
						scoreStrips = append(scoreStrips, ss)
					}
				}
				continue
			}
			//ss := parse(dbCOnn, season, PreSeason, week)
			//if ss.GameWeek.Games != nil {
			//	scoreStrips = append(scoreStrips, ss)
			//}

			ss := parse(dbCOnn, season, RegularSeason, week)
			if ss.GameWeek.Games != nil {
				scoreStrips = append(scoreStrips, ss)
			}

			if season != 2017 {
				ss = parse(dbCOnn, season, PostSeason, week)
				if ss.GameWeek.Games != nil {
					scoreStrips = append(scoreStrips, ss)
				}
			}

		}
	}
	return scoreStrips
}

func downloadAll(dbCOnn *DB) {
	for season := firstSeason; season <= lastSeason; season++ {
		for week := firstWeek; week <= lastWeek; week++ {
			if season == 2017 {
				scrape(dbCOnn, season, PreSeason, week)
				if week <= currentWeek {
					scrape(dbCOnn, season, RegularSeason, week)
				}
				continue
			}
			scrape(dbCOnn, season, PreSeason, week)
			scrape(dbCOnn, season, RegularSeason, week)
			scrape(dbCOnn, season, PostSeason, week)
		}
	}
}

func name(season int, phase string, week int) []byte {
	return []byte(fmt.Sprintf("%d_%s_%d", season, phase, week))
}

func scrape(db *DB, season int, phase string, week int) *ScoreStrip {
	var s ScoreStrip
	if !exists(db, season, phase, week) {
		download(db, season, phase, week)
	}
	return &s
}

func exists(db *DB, season int, phase string, week int) bool {
	var found bool
	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(nflRawXML))
		data := bucket.Get(name(season, phase, week))
		found = data != nil
		return nil
	})
	if err != nil {
		panic(err)
	}

	return found
}

func download(db *DB, season int, phase string, week int) {
	const scrapeURLFormat = "http://www.nfl.com/ajax/scorestrip?season=%d&seasonType=%s&week=%d"
	resp, err := http.Get(fmt.Sprintf(scrapeURLFormat, season, phase, week))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	write(db, season, phase, week, body)
}

func write(db *DB, season int, phase string, week int, data []byte) {
	if err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(nflRawXML))
		if err != nil {
			return err
		}
		fmt.Println("WRITING ", string(name(season, phase, week)))
		return bucket.Put(name(season, phase, week), data)
	}); err != nil {
		panic(err)
	}
}
