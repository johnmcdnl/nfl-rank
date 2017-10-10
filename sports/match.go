package sports

type Match struct {
	ID              string
	Name            string
	IsCompleted     bool
	HomeTeam        *Team
	AwayTeam        *Team
	WeightingFactor float64
}

func (m *Match) Winner() Result {
	if m.HomeTeam.Score > m.AwayTeam.Score {
		return HomeWin
	}
	if m.HomeTeam.Score < m.AwayTeam.Score {
		return AwayWin
	}
	return Draw
}
