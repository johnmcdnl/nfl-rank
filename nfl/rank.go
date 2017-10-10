package nfl

import (
	"fmt"
	"github.com/johnmcdnl/nfl-rank/sports"
	"github.com/johnmcdnl/elo"
	"sort"
)

const initialRating float64 = 1500

func CalculateELOForSeasons(seasons []*sports.Season) {
	for _, season := range seasons {
		for _, phase := range season.Phases {
			for _, gameWeek := range phase.GameWeeks {
				for _, match := range gameWeek.Matches {
					CalculateELO(match)
				}
			}
		}
	}
}

func validTeams() []string {
	var teams []string
	teams = append(teams, "patriots")
	teams = append(teams, "chiefs")
	teams = append(teams, "broncos")
	teams = append(teams, "packers")
	teams = append(teams, "seahawks")
	teams = append(teams, "steelers")
	teams = append(teams, "panthers")
	teams = append(teams, "falcons")
	teams = append(teams, "cardinals")
	teams = append(teams, "cowboys")
	teams = append(teams, "eagles")
	teams = append(teams, "lions")
	teams = append(teams, "bengals")
	teams = append(teams, "vikings")
	teams = append(teams, "colts")
	teams = append(teams, "dolphins")
	teams = append(teams, "texans")
	teams = append(teams, "ravens")
	teams = append(teams, "saints")
	teams = append(teams, "bills")
	teams = append(teams, "raiders")
	teams = append(teams, "redskins")
	teams = append(teams, "jets")
	teams = append(teams, "giants")
	teams = append(teams, "buccaneers")
	teams = append(teams, "rams")
	teams = append(teams, "titans")
	teams = append(teams, "chargers")
	teams = append(teams, "49ers")
	teams = append(teams, "bears")
	teams = append(teams, "jaguars")
	teams = append(teams, "browns")
}

func ReportELOs() {

	n := map[float64][]string{}
	var a []float64
	for k, v := range rankings {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(a)))
	for _, k := range a {
		for i, s := range n[k] {
			fmt.Printf("%d : %s %f\n", i, s, k)
		}
	}

}

var rankings = make(map[string]float64)

func getCurrentRank(name string) float64 {
	v, found := rankings[name]
	if !found {
		return initialRating
	}
	return v
}

func CalculateELO(match *sports.Match) {
	homeRank := getCurrentRank(match.HomeTeam.NickName)
	awayRank := getCurrentRank(match.AwayTeam.NickName)

	var result *elo.ELO
	var err error

	switch match.Winner() {
	default:
		panic("Unhandled exception")
	case sports.HomeWin:
		result, err = elo.New(homeRank, awayRank, 24, elo.Win, elo.Loose)
	case sports.AwayWin:
		result, err = elo.New(homeRank, awayRank, 24, elo.Loose, elo.Win)
	case sports.Draw:
		result, err = elo.New(homeRank, awayRank, 24, elo.Draw, elo.Draw)
	}

	if err != nil {
		panic(err)
	}

	rankings[match.HomeTeam.NickName] = result.RAN
	rankings[match.AwayTeam.NickName] = result.RBN

}
