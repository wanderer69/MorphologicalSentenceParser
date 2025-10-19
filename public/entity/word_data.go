package entity

type WordData struct {
	Rel      string
	Pos      string
	Feats    map[string]string
	Start    string
	Stop     string
	Text     string
	Lemma    string
	Id       string
	HeadID   string
	IdN      int
	SidN     int
	HeadIdN  int
	SheadIdN int
	StartN   int
	StopN    int
}
