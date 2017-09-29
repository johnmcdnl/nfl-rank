package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())
	nfl.GenerateELO(seasons)

}
