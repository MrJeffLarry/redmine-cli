package auth

type user struct {
	User userInfo `json:"user,omitempty"`
}

type userInfo struct {
	ID          int    `json:"id,omitempty"`
	Login       string `json:"login,omitempty"`
	Admin       bool   `json:"admin,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	TwofaScheme string `json:"twofa_scheme,omitempty"`
	ApiKey      string `json:"api_key,omitempty"`
}
