package resources

import (
	"net/http"
	"github.com/go-chi/render"
)

type Article struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

type ArticleResponse struct {
	*Article
	Elapsed int64 `json:"elapsed"`
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Elapsed = 10
	return nil
}

func NewArticleResponse(article *Article) *ArticleResponse {
	resp := &ArticleResponse{Article: article}
	return resp
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)
	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}