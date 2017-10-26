package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
	"fmt"
	"encoding/json"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())
	//
	ranker := nfl.NewRanker(seasons)

	ranker.PerformRanking()
	ranker.Report()

	j, _ := json.Marshal(ranker.ReportHistorical("eagles"))
	fmt.Sprintln(string(j))
}
