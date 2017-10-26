package resources

import (
	"net/http"
	"github.com/go-chi/render"
	"time"
)

type TeamRating struct {
	Name string `json:"name"`
	Ranks []struct {
		Time    time.Time `json:"time"`
		Ranking float64   `json:"ranking"`
	} `json:"ranks"`
}

type TeamRatingResponse struct {
	*TeamRating
	Elapsed int64 `json:"elapsed"`
}

func GetTeamRating(w http.ResponseWriter, r *http.Request) {
	teamRating := r.Context().Value("teamRating").(*TeamRating)
	if err := render.Render(w, r, NewTeamRatingResponse(teamRating)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func NewTeamRatingResponse(teamRating *TeamRating) *TeamRatingResponse {
	resp := &TeamRatingResponse{TeamRating: teamRating}
	return resp
}

func (rd *TeamRatingResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Elapsed = 10
	return nil
}
