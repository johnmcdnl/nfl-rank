package elo

import (
	"errors"
	"math"
)

type Result float64

const (
	Win   = Result(1)
	Draw  = Result(0.5)
	Loose = Result(0)
)

type ELO struct {
	RA  float64
	RB  float64
	K   float64
	EA  float64
	EB  float64
	SA  Result
	SB  Result
	RAN float64
	RBN float64
}

func New(rA, rB, k float64, sa, sb Result) (*ELO, error) {
	if float64(sa)+float64(sb) != 1 {
		return nil, errors.New("invalid result")
	}
	var e = &ELO{
		RA: rA,
		RB: rB,
		K:  k,
		SA: sa,
		SB: sb,
	}
	e.calculate()
	return e, nil
}

func (e *ELO) calculate() {
	qA := math.Pow(10, e.RA/400)
	qB := math.Pow(10, e.RB/400)

	e.EA = qA / (qA + qB)
	e.EB = qB / (qA + qB)

	e.RAN = e.RA + e.K*(float64(e.SA)-e.EA)
	e.RBN = e.RB + e.K*(float64(e.SB)-e.EB)
}
