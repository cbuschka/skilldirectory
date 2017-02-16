package model

/*
User represents
*/
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

/*

 */
type UserAccount struct {
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

func NewUser(login, password string) User {
	return User{
		Login:    login,
		Password: password,
	}
}

type AuthCredentials struct {
	Id     string `json:"client_id"`
	Secret string `json:"client_secret"`
	Code   string `json:"code"`
}

type TokenResponse struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}
