package model

type Group struct {
	Name        string
	Description string
	Members     []*Application
}

type GroupUpdate struct {
	Name        *string
	Description *string
	Members     []string
}
