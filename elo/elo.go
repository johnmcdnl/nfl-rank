package elo

import (
	"math"
)

type Result int8

const (
	A_WIN = Result(iota)
	B_WIN
	Draw
)

type ELO struct {
	result Result
	k      float64

	R_A float64
	R_B float64

	E_A float64
	E_B float64

	S_A float64
	S_B float64

	R_A_N float64
	R_B_N float64
}

func New(p1, p2 float64, result Result, k float64, homeAdvantage float64) *ELO {
	e := &ELO{
		R_A:    p1,
		R_B:    p2,
		result: result,
		k:      k,
	}

	switch result {
	case A_WIN:
		e.S_A = 1 - homeAdvantage
		e.S_B = 0 + homeAdvantage
	case B_WIN:
		e.S_A = 0 - homeAdvantage
		e.S_B = 1 + homeAdvantage
	case Draw:
		e.S_A = 0.5 - homeAdvantage
		e.S_B = 0.5 + homeAdvantage
	default:
		panic("Unhandled")
	}

	return e
}

func (e *ELO) Calculate() {
	e.E_A = 1 / (1 + math.Pow(10, (e.R_B-e.R_A)/400))
	e.E_B = 1 / (1 + math.Pow(10, (e.R_A-e.R_B)/400))
	e.R_A_N = e.R_A + e.S_A + e.k*(e.S_A-e.E_A)
	e.R_B_N = e.R_B + e.S_B +  e.k*(e.S_B-e.E_B)

}

func (e *ELO) NewP1() float64 {
	return e.R_A_N

}
func (e *ELO) NewP2() float64 {
	return e.R_B_N
}
