package sports

import "time"

type Match struct {
	ID          string
	Name        string
	IsCompleted bool
	HomeTeam    *Team
	AwayTeam    *Team
	Weight      float64
	Time        time.Time
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
