package model

type User struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
