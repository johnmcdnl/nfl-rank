package nfl

type ScoreStrip struct {
	GameWeek *GameWeek `xml:" ss,omitempty" json:"gameWeek,omitempty"`
}
type GameWeek struct {
	Games *Games `xml:" gms,omitempty" json:"games,omitempty"`
}

type Games struct {
	gd      string  `xml:" gd,attr"  json:",omitempty"`
	Phase   string  `xml:" t,attr"  json:",omitempty"`
	WeekNum string  `xml:" w,attr"  json:",omitempty"`
	Season  string  `xml:" y,attr"  json:",omitempty"`
	Games   []*Game `xml:" g,omitempty" json:"games,omitempty"`
}

type Game struct {
	WeekDay         string `xml:" d,attr"  json:"weekDay,omitempty"`
	EventID         string `xml:" eid,attr"  json:"eventID,omitempty"`
	ga              string `xml:" ga,attr"  json:",omitempty"`
	gsis            string `xml:" gsis,attr"  json:",omitempty"`
	GameType        string `xml:" gt,attr"  json:",omitempty"`
	Home            string `xml:" h,attr"  json:",omitempty"`
	HomeNickName    string `xml:" hnn,attr"  json:",omitempty"`
	HomeScore       int    `xml:" hs,attr"  json:",omitempty"`
	k               string `xml:" k,attr"  json:",omitempty"`
	p               string `xml:" p,attr"  json:",omitempty"`
	q               string `xml:" q,attr"  json:",omitempty"`
	rz              string `xml:" rz,attr"  json:",omitempty"`
	Time            string `xml:" t,attr"  json:",omitempty"`
	Visitor         string `xml:" v,attr"  json:",omitempty"`
	VisitorNickName string `xml:" vnn,attr"  json:",omitempty"`
	VisitorScore    int    `xml:" vs,attr"  json:",omitempty"`
}
