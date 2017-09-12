package sports

type GameWeek struct {
	ID      string
	Name    string
	Matches []*Match
}

func (week *GameWeek) NewMatch(match *Match) {
	week.Matches = append(week.Matches, match)
}
