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
