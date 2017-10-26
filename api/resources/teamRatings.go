package resources

import (
	"github.com/go-chi/render"
	"net/http"
	"github.com/johnmcdnl/nfl-rank/nfl"
)

type TeamRatings struct {
	TeamRatings []*TeamRating `json:"teams"`
}

func ListTeamRatings(w http.ResponseWriter, r *http.Request) {
	teamRatings := GenerateTeamRatings()
	//j, _ := json.Marshal(teamRatings)
	//w.Write(j)
	if err := render.RenderList(w, r, NewTeamRatingsListResponse(teamRatings.TeamRatings)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func NewTeamRatingsListResponse(teamRatings []*TeamRating) []render.Renderer {
	list := []render.Renderer{}
	for _, teamRating := range teamRatings {
		list = append(list, NewTeamRatingResponse(teamRating))
	}
	return list
}

func GenerateTeamRatings() *TeamRatings {

	seasons := nfl.Transform(nfl.ScrapeAll())
	ranker := nfl.NewRanker(seasons)
	ranker.PerformRanking()
	var rankings = TeamRatings{}
	for _, current := range nfl.CurrentTeams() {
		teamHistory := ranker.ReportHistorical(current)
		var tr TeamRating
		tr.Name = current

		for _, r := range teamHistory {
			var rank Rank
			rank.Time = r.Time
			rank.Ranking = r.RankingPoints
			tr.Ranks = append(tr.Ranks, rank)
		}

		rankings.TeamRatings = append(rankings.TeamRatings, &tr)
	}

	return &rankings
}
