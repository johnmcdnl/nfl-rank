package nfl

import (
	"github.com/johnmcdnl/nfl-rank/sports"
	"github.com/johnmcdnl/elo"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Ranks struct {
	Time          time.Time
	Team          *sports.Team
	RankingPoints float64
}

type Ranker struct {
	BaseRank float64
	K        float64
	HomeBias float64

	LowWeight     float64
	RegularWeight float64
	HighWeight    float64

	seasons []*sports.Season
	Teams   RankedTeams

	HistoricRanks []*Ranks

	ModelCorrect   int
	ModelIncorrect int
}

type RankedTeam struct {
	Team          *sports.Team
	WinLossRecord
	HomeBias      float64
	RankingPoints float64
	isNew         bool
}

type RankedTeams struct {
	Teams []*RankedTeam
}

func (teams *RankedTeams) GetTeam(team *sports.Team) *RankedTeam {
	for _, t := range teams.Teams {
		if t.Team.NickName == team.NickName {
			return t
		}
	}
	var rt = &RankedTeam{Team: team, isNew: true}
	teams.Teams = append(teams.Teams, rt)
	return rt
}

type WinLossRecord struct {
	HomeWins   int
	HomeLosses int
	HomeDraws  int
	AwayWins   int
	AwayLosses int
	AwayDraws  int
	Total      int
}

func NewRanker(seasons []*sports.Season) *Ranker {
	return &Ranker{
		BaseRank: 1500,
		K:        48,
		HomeBias: 30,

		HighWeight:    1.2,
		RegularWeight: 1.1,
		LowWeight:     0.3,

		seasons: seasons,
	}
}

func (r *Ranker) PerformRanking() {
	for _, s := range r.seasons {
		r.RankSeason(s)
	}
}

func (r *Ranker) RankSeason(s *sports.Season) {
	for _, phase := range s.Phases {
		r.RankPhase(phase)
	}
}

func (r *Ranker) RankPhase(p *sports.Phase) {
	var weight float64
	switch p.Name {
	default:
		panic("unhandled")
	case PreSeason:
		weight = r.LowWeight
	case RegularSeason:
		weight = r.RegularWeight
	case PostSeason:
		weight = r.HighWeight
	}
	for _, gw := range p.GameWeeks {
		r.RankGameWeek(gw, weight)
	}
}

func (r *Ranker) RankGameWeek(gw *sports.GameWeek, weighting float64) {
	for _, m := range gw.Matches {
		m.Weight = weighting
		r.RankMatch(m)
	}
}

func (r *Ranker) RankMatch(m *sports.Match) {
	if !m.IsCompleted {
		return
	}
	r.CalculateELO(m)
	r.CalculateWinLossForMatch(m)

	r.HistoricRanks = append(r.HistoricRanks, &Ranks{Team: m.HomeTeam, RankingPoints: r.GetTeam(m.HomeTeam).RankingPoints, Time: m.Time}, )
	r.HistoricRanks = append(r.HistoricRanks, &Ranks{Team: m.AwayTeam, RankingPoints: r.GetTeam(m.AwayTeam).RankingPoints, Time: m.Time}, )
}

func (r *Ranker) CalculateWinLossForMatch(m *sports.Match) {
	if m.Winner() == sports.HomeWin {
		r.GetTeam(m.HomeTeam).WinLossRecord.HomeWins++
		r.GetTeam(m.AwayTeam).WinLossRecord.AwayWins--
	}

	if m.Winner() == sports.AwayWin {
		r.GetTeam(m.HomeTeam).WinLossRecord.HomeLosses++
		r.GetTeam(m.AwayTeam).WinLossRecord.AwayWins++
	}

	if m.Winner() == sports.Draw {
		r.GetTeam(m.HomeTeam).WinLossRecord.HomeDraws++
		r.GetTeam(m.AwayTeam).WinLossRecord.AwayDraws++
	}
}

func (r *Ranker) GetTeam(t *sports.Team) *RankedTeam {
	team := r.Teams.GetTeam(t)
	if team.isNew {
		team.RankingPoints = r.BaseRank
		team.isNew = false
	}
	return team
}

func (r *Ranker) CalculateELO(m *sports.Match) {
	home := r.GetTeam(m.HomeTeam)
	away := r.GetTeam(m.AwayTeam)

	hRating := home.RankingPoints
	aRating := away.RankingPoints

	hBiased := home.RankingPoints + r.HomeBias
	aBiased := away.RankingPoints - r.HomeBias

	kWeight := r.K * m.Weight

	var result *elo.ELO
	var err error

	switch m.Winner() {
	default:
		panic("Unhandled exception")
	case sports.HomeWin:
		result, err = elo.New(hBiased, aBiased, kWeight, elo.Win, elo.Loose)
	case sports.AwayWin:
		result, err = elo.New(hBiased, aBiased, kWeight, elo.Loose, elo.Win)
	case sports.Draw:
		result, err = elo.New(hBiased, aBiased, kWeight, elo.Draw, elo.Draw)
	}
	if err != nil {
		panic(err)
	}

	home.RankingPoints = hRating + (result.RAN - hRating) - r.HomeBias
	away.RankingPoints = aRating + (result.RBN - aRating) + r.HomeBias

	switch m.Winner() {
	default:
		panic("Unhandled exception")
	case sports.HomeWin:
		if result.EA > result.EB {
			r.ModelCorrect++
		} else {
			r.ModelIncorrect++
		}
	case sports.AwayWin:
		if result.EB > result.EA {
			r.ModelCorrect++
		} else {
			r.ModelIncorrect++
		}
	case sports.Draw:
		if result.EA == result.EB {
			r.ModelCorrect++
		} else {
			r.ModelIncorrect++
		}
	}

}

func (r *Ranker) Report() {
	r.Sort()
	var cRank = 1
	for _, t := range r.Teams.Teams {
		if !currentTeam(t.Team.NickName) {
			continue
		}
		fmt.Println(cRank, t.RankingPoints, t.Team.NickName)
		cRank++
	}
}

func (r *Ranker) ReportHistorical(nickname string) {

	for _, hr := range r.HistoricRanks {
		if strings.ToLower(hr.Team.NickName) == strings.ToLower(nickname) {
			fmt.Println(hr.Time.Format(time.RFC3339), hr.RankingPoints)
		}
	}
}
func (wlr WinLossRecord) Percentage() float64 {
	wins := wlr.HomeWins + wlr.AwayWins + ((wlr.HomeDraws + wlr.AwayDraws) / 2)
	losses := wlr.HomeLosses + wlr.AwayLosses

	return float64(wins) / (float64(wins) + float64(losses))
}

func (r *Ranker) Sort() {
	sort.Slice(r.Teams.Teams, func(i, j int) bool { return r.Teams.Teams[i].RankingPoints > r.Teams.Teams[j].RankingPoints })
}

func (r *Ranker) Accuracy() float64 {
	return float64(r.ModelCorrect) / (float64(r.ModelCorrect + r.ModelIncorrect))
}

func currentTeam(name string) bool {
	switch strings.ToLower(name) {
	default:
		return false
	case "patriots", "dolphins", "jets", "bills":
		return true
	case "chiefs", "broncos", "raiders", "chargers":
		return true
	case "steelers", "ravens", "browns", "bengals":
		return true
	case "texans", "colts", "titans", "jaguars":
		return true
	case "cowboys", "eagles", "giants", "redskins":
		return true
	case "seahawks", "cardinals", "rams", "49ers":
		return true
	case "packers", "lions", "vikings", "bears":
		return true
	case "panthers", "falcons", "saints", "buccaneers":
		return true
	}
}
