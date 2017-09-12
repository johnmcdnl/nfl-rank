package sports

type Seasons struct {
	Seasons []*Season
}

func (seasons *Seasons) GetSeasonNamed(name string) *Season {
	for _, season := range seasons.Seasons {
		if season.Name == name {
			return season
		}
	}

	var season = &Season{}
	season.Name = name
	seasons.Seasons = append(seasons.Seasons, season)
	return seasons.GetSeasonNamed(name)
}

type Season struct {
	ID     string
	Name   string
	Phases []*Phase
}

func (season *Season) SortPhases(phases []*Phase, order []string) {
	var sorted []*Phase

	for _, o := range order {
		sorted = append(sorted, season.GetPhaseNamed(o))
	}

	season.Phases = sorted
}

func (season *Season) GetPhaseNamed(name string) *Phase {
	for _, phase := range season.Phases {
		if phase.Name == name {
			return phase
		}
	}

	var phase = &Phase{}
	phase.Name = name
	season.Phases = append(season.Phases, phase)
	return season.GetPhaseNamed(name)
}
