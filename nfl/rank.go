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
		for _, s := range n[k] {
			fmt.Printf("%s %f\n", s, k)
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
