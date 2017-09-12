package sports

type Phase struct {
	ID        string
	Name      string
	GameWeeks []*GameWeek
}

func (phase *Phase) GetGameWeekNamed(name string) *GameWeek {
	for _, gameWeek := range phase.GameWeeks {
		if gameWeek.Name == name {
			return gameWeek
		}
	}
	var gameWeek = &GameWeek{}
	gameWeek.Name = name
	phase.GameWeeks = append(phase.GameWeeks, gameWeek)
	return gameWeek
}
