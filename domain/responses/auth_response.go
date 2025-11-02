package responses

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	Sub       string `json:"sub"`
	Role      string `json:"role"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsActive  bool   `json:"isActive"`
}

type MeResponse struct {
	Authenticated bool     `json:"authenticated"`
	User          UserInfo `json:"user,omitempty"`
}
