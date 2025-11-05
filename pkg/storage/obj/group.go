package obj

type Group struct {
	CosmosObj
	Name         string `gorm:"uniqueIndex"`
	Description  string
	Applications []*Application `gorm:"many2many:group_applications;"`
}
