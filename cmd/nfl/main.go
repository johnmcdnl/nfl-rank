package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
	"github.com/johnmcdnl/nfl-rank/sports"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())

	j, _ := json.Marshal(seasons)
	ioutil.WriteFile("season.json", j, os.ModePerm)

}

func preSeason(s *sports.Season) *sports.Phase {
	return s.GetPhaseNamed(nfl.PreSeason)
}

func regularSeason(s *sports.Season) *sports.Phase {
	return s.GetPhaseNamed(nfl.RegularSeason)
}

func postSeason(s *sports.Season) *sports.Phase {
	return s.GetPhaseNamed(nfl.PostSeason)
}

func printPhase(p *sports.Phase) {
	for i := nfl.FirstWeek; i <= nfl.LastWeek; i++ {
		gw := p.GetGameWeekNamed(fmt.Sprint(i))
		printGameWeek(p, gw)
	}
}

func printGameWeek(p *sports.Phase, gw *sports.GameWeek) {
	for _, g := range gw.Matches {
		fmt.Println(p.Name, gw.Name, g.HomeTeam, g.AwayTeam)
	}
}
