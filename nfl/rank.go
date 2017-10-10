package nfl

import (
	"fmt"
	"github.com/johnmcdnl/nfl-rank/sports"
	"github.com/johnmcdnl/elo"
	"sort"
	"strings"
	"sync"
)

const initialRating float64 = 1500
const kFactor float64 = 60

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
						match.WeightingFactor = 0.3
					case RegularSeason:
						match.WeightingFactor = 1.1
					case PostSeason:
						match.WeightingFactor = 1.2
					}
					CalculateELO(match)
				}
			}
		}
	}
}

func FuturePredictions(seasons []*sports.Season) {
	nextWeekFound := false
	for _, season := range seasons {
		for _, phase := range season.Phases {
			for _, gameWeek := range phase.GameWeeks {
				for _, match := range gameWeek.Matches {
					if match.IsCompleted {
						continue
					}
					switch phase.Name {
					case PreSeason:
						match.WeightingFactor = 0.3
					case RegularSeason:
						match.WeightingFactor = 1.1
					case PostSeason:
						match.WeightingFactor = 1.2
					}
					once.Do(func(){
						fmt.Printf("%s - %s - %s \n", season.Name, phase.Name, gameWeek.Name)
					})
					PredictWinner(match)
					nextWeekFound = true
				}
				once = sync.Once{}
				//if nextWeekFound {
				//	return
				//}
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
	homeRank := getCurrentRank(match.HomeTeam.NickName)
	awayRank := getCurrentRank(match.AwayTeam.NickName)

	var result *elo.ELO
	var err error

	switch match.Winner() {
	default:
		panic("Unhandled exception")
	case sports.HomeWin:
		result, err = elo.New(homeRank, awayRank, kFactor*match.WeightingFactor, elo.Win, elo.Loose)
	case sports.AwayWin:
		result, err = elo.New(homeRank, awayRank, kFactor*match.WeightingFactor, elo.Loose, elo.Win)
	case sports.Draw:
		result, err = elo.New(homeRank, awayRank, kFactor*match.WeightingFactor, elo.Draw, elo.Draw)
	}

	if err != nil {
		panic(err)
	}

	rankings[match.HomeTeam.NickName] = result.RAN
	rankings[match.AwayTeam.NickName] = result.RBN

}

func PredictWinner(match *sports.Match) {
	homeRank := getCurrentRank(match.HomeTeam.NickName)
	awayRank := getCurrentRank(match.AwayTeam.NickName)

	result, err := elo.New(homeRank, awayRank, kFactor*match.WeightingFactor, elo.Win, elo.Loose)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%2.0f%s \t %s %2.2f @ %2.2f %s \t %2.0f%s \n",
		result.EB*100, "%",
		PadLeft(match.AwayTeam.NickName, " ", 10),

		100/(result.EB*100),
		100/(result.EA*100),

		PadRight(match.HomeTeam.NickName, " ", 10),
		result.EA*100, "%",
	)

}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) >= lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) >= lenght {
			return str[0:lenght]
		}
	}
}
