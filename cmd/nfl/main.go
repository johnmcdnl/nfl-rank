package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())
	//nfl.CalculateELOForSeasons(seasons)
	//nfl.ReportELOs()
	//fmt.Println()
	//fmt.Println()
	//nfl.FuturePredictions(seasons)

	ranker := nfl.NewRanker(seasons)
	ranker.PerformRanking()
	ranker.Report()

	ranker.ReportHistorical("eagles")
	//nfl.EstimationBestValues(seasons)
}
