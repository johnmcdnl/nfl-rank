package sports

type GameWeek struct {
	ID       string
	Name     string
	GameType string
	Matches  []*Match
}

func (week *GameWeek) NewMatch(match *Match) {
	week.Matches = append(week.Matches, match)
}
