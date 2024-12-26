package web_pkg

/*
POST to set team state
*/
type SetTeam struct {
	Slug          string `json:"slug,omitempty"`
	ReviewEnabled bool   `json:"review_enabled,omitempty"`
	ReviewOrder   int    `json:"review_order,omitempty"`
}

/*
POST to set repo state
*/
type setRepo struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
