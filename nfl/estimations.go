package nfl

import (
	"github.com/johnmcdnl/nfl-rank/sports"
	"fmt"
)

func EstimationBestValues(seasons []*sports.Season) {
	var bestEstimation float64
	var bestEstimationK int
	var bestEstimationHome int

	for k := 0; k <= 200; k++ {
		for h := 0; h <= 200; h++ {
			r := NewRanker(seasons)
			r.K=float64(k)
			r.HomeBias= float64(h)
			r.PerformRanking()
			a := r.Accuracy()
			if a > bestEstimation {
				bestEstimation = a
				bestEstimationK = k
				bestEstimationHome = h
				fmt.Println(k, h, a)
			}

		}
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("BEST")
	fmt.Println(bestEstimationK, bestEstimationHome, bestEstimation)
}
