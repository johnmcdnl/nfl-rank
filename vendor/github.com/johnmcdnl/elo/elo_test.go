package elo

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestNew(t *testing.T) {
	type args struct {
		rA float64
		rB float64
		k  float64
		sa Result
		sb Result
	}
	tests := []struct {
		name    string
		args    args
		want    *ELO
		wantErr bool
	}{
		{"Fav Win", args{rA: 1700, rB: 1500, k: 20, sa: Win, sb: Loose}, &ELO{RA: 1700, RB: 1500, K: 20, SA: Win, SB: Loose, EA: 0.75974, EB: 0.24025, RAN: 1704.80506, RBN: 1495.19493}, false},
		{"Fav Draw", args{rA: 1700, rB: 1500, k: 20, sa: Draw, sb: Draw}, &ELO{RA: 1700, RB: 1500, K: 20, SA: Draw, SB: Draw, EA: 0.75974, EB: 0.24025, RAN: 1694.80506, RBN: 1505.19493}, false},
		{"Fav Loose", args{rA: 1700, rB: 1500, k: 20, sa: Loose, sb: Win}, &ELO{RA: 1700, RB: 1500, K: 20, SA: Loose, SB: Win, EA: 0.75974, EB: 0.24025, RAN: 1684.80506, RBN: 1515.19493}, false},

		{"Underdog Win", args{rA: 1500, rB: 1700, k: 20, sa: Win, sb: Loose}, &ELO{RA: 1500, RB: 1700, K: 20, SA: Win, SB: Loose, EA: 0.24025, EB: 0.75974, RAN: 1515.19493, RBN: 1684.80506}, false},
		{"Underdog Draw", args{rA: 1500, rB: 1700, k: 20, sa: Draw, sb: Draw}, &ELO{RA: 1500, RB: 1700, K: 20, SA: Draw, SB: Draw, EA: 0.24025, EB: 0.75974, RAN: 1505.19493, RBN: 1694.80506}, false},
		{"Underdog Loose", args{rA: 1500, rB: 1700, k: 20, sa: Loose, sb: Win}, &ELO{RA: 1500, RB: 1700, K: 20, SA: Loose, SB: Win, EA: 0.24025, EB: 0.75974, RAN: 1495.19493, RBN: 1704.80506}, false},

		{"Equal Rank Win", args{rA: 1500, rB: 1500, k: 20, sa: Win, sb: Loose}, &ELO{RA: 1500, RB: 1500, K: 20, SA: Win, SB: Loose, EA: 0.5, EB: 0.5, RAN: 1510, RBN: 1490}, false},
		{"Equal Rank Draw", args{rA: 1500, rB: 1500, k: 20, sa: Draw, sb: Draw}, &ELO{RA: 1500, RB: 1500, K: 20, SA: Draw, SB: Draw, EA: 0.5, EB: 0.5, RAN: 1500, RBN: 1500}, false},
		{"Equal Rank Loose", args{rA: 1500, rB: 1500, k: 20, sa: Loose, sb: Win}, &ELO{RA: 1500, RB: 1500, K: 20, SA: Loose, SB: Win, EA: 0.5, EB: 0.5, RAN: 1490, RBN: 1510}, false},

		{"2 Winners", args{rA: 1500, rB: 1500, k: 20, sa: Win, sb: Win}, nil, true},
		{"2 Looser", args{rA: 1500, rB: 1500, k: 20, sa: Loose, sb: Loose}, nil, true},
		{"Win & Draw", args{rA: 1500, rB: 1500, k: 20, sa: Win, sb: Draw}, nil, true},
		{"Loose & Draw", args{rA: 1500, rB: 1500, k: 20, sa: Loose, sb: Draw}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.rA, tt.args.rB, tt.args.k, tt.args.sa, tt.args.sb)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)

			const epsilon = 0.0001

			assert.Equal(t, tt.want.SA, got.SA, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "SA", float64(tt.want.SA), float64(got.SA)))
			assert.Equal(t, tt.want.SB, got.SB, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "SA", float64(tt.want.SB), float64(got.SB)))

			assert.InEpsilon(t, tt.want.K, got.K, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "K", float64(tt.want.K), float64(got.K)))

			assert.InEpsilon(t, tt.want.RA, got.RA, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "RA", float64(tt.want.RA), float64(got.RA)))
			assert.InEpsilon(t, tt.want.RB, got.RB, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "RB", float64(tt.want.RB), float64(got.RB)))

			assert.InEpsilon(t, tt.want.EA, got.EA, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "EA", float64(tt.want.EA), float64(got.EA)))
			assert.InEpsilon(t, tt.want.EB, got.EB, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "EB", float64(tt.want.EB), float64(got.EB)))

			assert.InEpsilon(t, tt.want.RAN, got.RAN, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "RAN", float64(tt.want.RAN), float64(got.RAN)))
			assert.InEpsilon(t, tt.want.RBN, got.RBN, epsilon, fmt.Sprintf("%s - Expected: %2g ::: %2g Actual", "RBN", float64(tt.want.RBN), float64(got.RBN)))
		})
	}
}
