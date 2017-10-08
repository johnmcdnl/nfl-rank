package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
	"fmt"
)

func main() {
	scoreStrips := nfl.ScrapeAll()
	for _, s := range scoreStrips {
		fmt.Println(s.GameWeek.Games)
	}
}
