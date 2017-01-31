package model

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserAccount struct {
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}
