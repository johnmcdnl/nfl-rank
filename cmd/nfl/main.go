package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
)

func main() {
	nfl.CalculateELOForSeasons(nfl.Transform(nfl.ScrapeAll()))
	nfl.ReportELOs()
}
