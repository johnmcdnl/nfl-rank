package nfl

import (
	"github.com/johnmcdnl/nfl-rank/sports"
	"fmt"
	"time"
)

func Transform(scoreStrips []*ScoreStrip) []*sports.Season {

	var seasons = new(sports.Seasons)

	for _, scoreStrip := range scoreStrips {
		season := seasons.GetSeasonNamed(scoreStrip.GameWeek.Games.Season)

		gms := scoreStrip.GameWeek.Games

		if gms == nil {
			return nil
		}

		season.Name = fmt.Sprint(gms.Season)

		for _, game := range gms.Games {
			addGameToSeason(season, game, gms.WeekNum)
		}

	}

	for _, s := range seasons.Seasons {
		s.SortPhases(s.Phases, []string{PreSeason, RegularSeason, PostSeason})
	}

	return seasons.Seasons
}

func addGameToSeason(season *sports.Season, game *Game, weekNum string) {

	switch game.GameType {
	case "PRE":
		addGameToPhase(season, PreSeason, game, weekNum)
	case "REG":
		addGameToPhase(season, RegularSeason, game, weekNum)
	case "WC":
		addGameToPhase(season, PostSeason, game, weekNum)
	case "DIV":
		addGameToPhase(season, PostSeason, game, weekNum)
	case "CON":
		addGameToPhase(season, PostSeason, game, weekNum)
	case "SB":
		addGameToPhase(season, PostSeason, game, weekNum)
	default:
		panic(game.GameType)
	}

}

func addGameToPhase(season *sports.Season, phaseName string, game *Game, weekNum string) {

	eventTime, err := time.Parse("20060102", string([]byte(game.EventID)[0:8]))
	if err != nil {
		panic(err)
	}
	isCompleted := eventTime.Before(time.Now())

	season.GetPhaseNamed(phaseName).GetGameWeekNamed(weekNum, game.GameType).NewMatch(&sports.Match{
		Name:        game.EventID,
		IsCompleted: isCompleted,
		HomeTeam: &sports.Team{
			Name:     game.Home,
			NickName: game.HomeNickName,
			Score:    game.HomeScore,
		},
		AwayTeam: &sports.Team{
			Name:     game.Visitor,
			NickName: game.VisitorNickName,
			Score:    game.VisitorScore,
		},
	})
}
