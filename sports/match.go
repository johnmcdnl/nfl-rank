package sports

type Match struct {
	ID       string
	Name     string
	HomeTeam *Team
	AwayTeam *Team
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
