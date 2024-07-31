package models

type VerifyPass struct {
	UserPassword  string `json:"userpassword"`
	InputPassword string `json:"inputpassword"`
}
