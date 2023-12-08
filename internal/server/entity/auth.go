package entity

type UserCredentials struct {
	Login      string
	Password   string
	Permission int
}

type ResourceInfo struct {
	Name   string
	Method string
}
