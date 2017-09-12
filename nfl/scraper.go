package nfl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"encoding/xml"
	"net/http"
)

const (
	MaxWorker = 50
	MaxQueue  = 25
)

const scrapeURLFormat = "http://www.nfl.com/ajax/scorestrip?season=%d&seasonType=%s&week=%d"

const (
	firstSeason = 1900
	lastSeason  = 2016

	FirstWeek = 0
	LastWeek  = 30

	PreSeason     = "PRE"
	RegularSeason = "REG"
	PostSeason    = "POST"

	scrapeDir = "./data/nfl"
)

func ScrapeAll() []*ScoreStrip {
	urls := generateURLS()
	downloadURLAndWriteToDisk(urls)
	scoreStrips := scrapeFromDisk()
	return scoreStrips
}

func generateURLS() []*url {
	var urls []*url
	for year := firstSeason; year <= lastSeason; year++ {
		for week := FirstWeek; week <= LastWeek; week++ {
			urls = append(urls, generateURL(year, PreSeason, week))
			urls = append(urls, generateURL(year, RegularSeason, week))
			urls = append(urls, generateURL(year, PostSeason, week))
		}
	}
	return urls
}

func scrapeFromDisk() []*ScoreStrip {
	os.MkdirAll(scrapeDir, os.ModePerm)
	var scoreStrips []*ScoreStrip
	err := filepath.Walk(scrapeDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		var scoreStrip ScoreStrip
		scoreStrip.GameWeek = new(GameWeek)

		err = xml.Unmarshal(contents, scoreStrip.GameWeek)
		if err != nil {
			os.RemoveAll(path)
			fmt.Println(path)
			panic(err)
		}
		if scoreStrip.GameWeek == nil || scoreStrip.GameWeek.Games == nil {
			return nil
		}
		scoreStrips = append(scoreStrips, &scoreStrip)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return scoreStrips
}

func generateURL(year int, phase string, week int) *url {
	dir := fmt.Sprintf("%s/%d/%s", scrapeDir, year, phase)

	fileName := fmt.Sprintf("%s/%02d.xml", dir, week)
	f, err := os.Open(fileName)
	defer f.Close()
	if err == nil {
		// Already exists - don't do it again
		return nil
	}

	var url url
	url.path = fmt.Sprintf(scrapeURLFormat, year, phase, week)
	url.dir = dir
	url.filename = fileName
	return &url
}

func downloadURL(url *url) {
	if url == nil {
		return
	}

	fmt.Println(url.path)
	resp, err := http.Get(url.path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(url.dir, os.ModePerm)
	err = ioutil.WriteFile(url.filename, contents, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func downloadURLAndWriteToDisk(targets []*url) {
	jobs := make(chan *url, MaxQueue)

	for i := 1; i <= MaxWorker; i++ {
		go func(i int) {
			for j := range jobs {
				downloadURL(j)
			}
		}(i)
	}

	for _, t := range targets {
		jobs <- t
	}
}

type url struct {
	path     string
	dir      string
	filename string
}
