package main

import (
	"github.com/johnmcdnl/nfl-rank/nfl"
	"fmt"
)

func main() {
	seasons := nfl.Transform(nfl.ScrapeAll())
	//
	ranker := nfl.NewRanker(seasons)
	//ranker.K = 48
	//ranker.HomeBias = 30
	ranker.PerformRanking()
	ranker.Report()
	fmt.Println(ranker.Accuracy())
	//var m sports.Match
	//m.HomeTeam = new(sports.Team)
	//m.AwayTeam = new(sports.Team)
	//m.HomeTeam.Name = "h"
	//m.HomeTeam.Score = 12
	//m.AwayTeam.Name = "a"
	//m.AwayTeam.Score = 14
	//m.Weight = 1
	//ranker.GetTeam(m.HomeTeam).RankingPoints = 1600
	//ranker.GetTeam(m.AwayTeam).RankingPoints = 1500
	//ranker.HomeBias = 0
	//
	//ranker.CalculateELO(&m)
	//ranker.ReportHistorical("eagles")
	nfl.EstimationBestValues(seasons)
}
