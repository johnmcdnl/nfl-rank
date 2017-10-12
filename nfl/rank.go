package nfl

import (
	"fmt"
	"github.com/johnmcdnl/nfl-rank/sports"
	"github.com/johnmcdnl/elo"
	"sort"
	"strings"
	"sync"
	"github.com/johnmcdnl/statistics"
)

const initialRating float64 = 1500

var KFactor float64 = 75
var HomeAdvantage float64 = 40

var homeWins, awayWins, draws int

func CalculateELOForSeasons(seasons []*sports.Season) {
	for _, season := range seasons {
		for _, phase := range season.Phases {
			for _, gameWeek := range phase.GameWeeks {
				for _, match := range gameWeek.Matches {
					if !match.IsCompleted {
						continue
					}
					switch phase.Name {
					case PreSeason:
						match.Weight = 0.3
					case RegularSeason:
						match.Weight = 1.1
					case PostSeason:
						match.Weight = 1.2
					}
					CalculateELO(match)
				}
			}
		}
		condenseRankings()
	}
}

var modelCorrect, modelInCorrect int

func Accuracy() float64 {
	return float64(modelCorrect) / (float64(modelCorrect) + float64(modelInCorrect))
}

func Reset() {
	modelCorrect = 0
	modelInCorrect = 0
	rankings = make(map[string]float64)
}

func FuturePredictions(seasons []*sports.Season) {
	futurePredictions := 0
	for _, season := range seasons {
		for _, phase := range season.Phases {
			for _, gameWeek := range phase.GameWeeks {
				var isFutureWeek = false
				for _, match := range gameWeek.Matches {
					if match.IsCompleted {
						continue
					}
					isFutureWeek = true
					switch phase.Name {
					case PreSeason:
						match.Weight = 0.3
					case RegularSeason:
						match.Weight = 1.1
					case PostSeason:
						match.Weight = 1.2
					}
					once.Do(func() {
						fmt.Printf("%s - %s - %s \n", season.Name, phase.Name, gameWeek.Name)
					})
					PredictWinner(match)
				}
				once = sync.Once{}
				if isFutureWeek {
					futurePredictions++
					return
				}
				if futurePredictions > 3 {
					return
				}
			}
		}
	}
}

var once sync.Once

func currentTeam(name string) bool {
	switch strings.ToLower(name) {
	default:
		return false
	case "patriots", "dolphins", "jets", "bills":
		return true
	case "chiefs", "broncos", "raiders", "chargers":
		return true
	case "steelers", "ravens", "browns", "bengals":
		return true
	case "texans", "colts", "titans", "jaguars":
		return true
	case "cowboys", "eagles", "giants", "redskins":
		return true
	case "seahawks", "cardinals", "rams", "49ers":
		return true
	case "packers", "lions", "vikings", "bears":
		return true
	case "panthers", "falcons", "saints", "buccaneers":
		return true
	}
}

func ReportELOs() {

	n := map[float64][]string{}
	var sorted []float64
	for k, v := range rankings {
		n[v] = append(n[v], k)
	}
	for k := range n {
		sorted = append(sorted, k)
	}

	var currentRank = 1
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))
	for _, r := range sorted {
		for _, n := range n[r] {
			if !currentTeam(n) {
				continue
			}
			fmt.Printf("%2.0f : %.2f : %s \n", float64(currentRank), r, n)
			currentRank++
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
	homeRank := getCurrentRank(match.HomeTeam.NickName) + HomeAdvantage
	awayRank := getCurrentRank(match.AwayTeam.NickName) - HomeAdvantage

	var result *elo.ELO
	var err error

	switch match.Winner() {
	default:
		panic("Unhandled exception")
	case sports.HomeWin:
		homeWins++
		result, err = elo.New(homeRank, awayRank, KFactor*match.Weight, elo.Win, elo.Loose)
	case sports.AwayWin:
		awayWins++
		result, err = elo.New(homeRank, awayRank, KFactor*match.Weight, elo.Loose, elo.Win)
	case sports.Draw:
		draws++
		result, err = elo.New(homeRank, awayRank, KFactor*match.Weight, elo.Draw, elo.Draw)
	}

	if err != nil {
		panic(err)
	}

	rankings[match.HomeTeam.NickName] = result.RAN
	rankings[match.AwayTeam.NickName] = result.RBN

	if result.EA > result.EB && match.Winner() == sports.HomeWin {
		modelCorrect++
	} else if result.EB > result.EA && match.Winner() == sports.AwayWin {
		modelCorrect++
	} else {
		modelInCorrect++
	}

}

func PredictWinner(match *sports.Match) {
	homeRank := getCurrentRank(match.HomeTeam.NickName) + HomeAdvantage
	awayRank := getCurrentRank(match.AwayTeam.NickName) - HomeAdvantage

	result, err := elo.New(homeRank, awayRank, KFactor*match.Weight, elo.Win, elo.Loose)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%2.0f%s \t %s %2.2f @ %2.2f %s \t %2.0f%s \n",
		result.EB*100, "%",
		PadLeft(match.AwayTeam.NickName, " ", 17),

		100/(result.EB*100),
		100/(result.EA*100),

		PadRight(match.HomeTeam.NickName, " ", 17),
		result.EA*100, "%",
	)

}

func PadRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) >= length {
			return str[0:length]
		}
	}
}

func PadLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str[0:length]
		}
	}
}

func condenseRankings() {
	var endOfSeason []float64
	for n, r := range rankings {
		if !currentTeam(n) {
			continue
		}
		endOfSeason = append(endOfSeason, r)
	}
	mean := statistics.Mean(endOfSeason)
	stdDev := statistics.StandardDeviationPopulation(endOfSeason)
	lowerBound := mean - stdDev
	upperBound := mean + stdDev


	var balancePoints float64

	var sumRankings float64
	for i, eosR := range endOfSeason {
		if eosR > upperBound {
			balancePoints += eosR - upperBound
		}
		if eosR < lowerBound {
			balancePoints += lowerBound - eosR
		}
		sumRankings+=float64(i+1)
	}


	n := map[float64][]string{}
	var sorted []float64
	for k, v := range rankings {
		n[v] = append(n[v], k)
	}
	for k := range n {
		sorted = append(sorted, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))
	for _, r := range sorted {
		for _, n := range n[r] {
			if !currentTeam(n) {
				continue
			}
				fmt.Println(n, r )

		}
	}



	fmt.Println(mean, stdDev, lowerBound, upperBound, balancePoints, sumRankings)
}
