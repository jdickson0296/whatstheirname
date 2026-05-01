package handlers

type MediaResult struct {
	Title   string  `json:"title"`
	Year    string     `json:"year"`
	Actors  string  `json:"actors"`
	Type    string  `json:"type"`
	Genre   string  `json:"genre"`
	Summary string  `json:"summary"`
	Rating  float64 `json:"rating"`
}
