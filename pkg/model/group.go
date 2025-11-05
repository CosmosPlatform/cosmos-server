package model

type Group struct {
	Name        string
	Description string
	Members     []*Application
}
