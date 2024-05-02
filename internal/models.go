package models

type App struct {
	UUID   string
	Name   string
	Secret string
}

type User struct {
	UUID     string
	Login    string
	Email    string
	PassHash []byte
}
