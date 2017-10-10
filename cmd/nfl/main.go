package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
	"fmt"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())
	nfl.CalculateELOForSeasons(seasons)
	nfl.ReportELOs()
	fmt.Println()
	fmt.Println()
	nfl.FuturePredictions(seasons)
}
