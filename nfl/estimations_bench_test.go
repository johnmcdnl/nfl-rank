package nfl

import "testing"

func BenchmarkEstimationBestValues(b *testing.B) {
	seasons := Transform(ScrapeAll())
	for i := 0; i <= b.N; i++ {
		EstimationBestValues(seasons)
	}

}
