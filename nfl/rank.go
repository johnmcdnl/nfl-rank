package nfl

import (
	"github.com/johnmcdnl/nfl-rank/sports"
	"fmt"
	"github.com/johnmcdnl/elo"
)

const initialRating float64 = 1500

var homeWin, awayWin, draw int

func homeAdvantage() float64 {
	total := homeWin + awayWin + draw
	if total == 0 {
		total = 1
	}
	homePercentage := (float64(homeWin) + float64(draw)/float64(2)) / float64(total)
	return 0.5 - homePercentage
}

var rankings = make(map[string]float64)

func GenerateELO(seasons []*sports.Season) {
	for _, season := range seasons {
		for _, phase := range season.Phases {
			for _, gameWeek := range phase.GameWeeks {
				for _, match := range gameWeek.Matches {
					GenerateELOForMatch(match, getKFactor(gameWeek.GameType), homeAdvantage())
					switch match.Winner() {
					case sports.HomeWin:
						homeWin++
					case sports.AwayWin:
						awayWin++
					case sports.Draw:
						draw++
					}

				}
			}
		}
	}

	for _, v := range rankings {
		fmt.Println(v)
	}
}

func getKFactor(matchType string) float64 {
	switch matchType {
	case "PRE":
		return 5
	case "REG":
		return 15
	case "WC":
		return 20
	case "DIV":
		return 20
	case "CON":
		return 25
	case "SB":
		return 35
	default:
		panic(matchType)
	}
}

func GenerateELOForMatch(match *sports.Match, kFactor float64, homeAdvantage float64) {

	homeStartRank, homeFound := rankings[match.HomeTeam.NickName]
	if !homeFound {
		homeStartRank = initialRating
	}

	awayStartRank, awayFound := rankings[match.HomeTeam.NickName]
	if !awayFound {
		awayStartRank = initialRating
	}

	e, _ := elo.New(homeStartRank, awayStartRank, kFactor, elo.Win, elo.Loose)

	rankings[match.HomeTeam.NickName] = e.RAN
	rankings[match.AwayTeam.NickName] = e.RBN
}
