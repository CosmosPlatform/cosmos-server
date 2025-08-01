package obj

type Team struct {
	CosmosObj
	Name        string `gorm:"uniqueIndex"`
	Description string
}
