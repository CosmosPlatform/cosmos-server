package obj

type Application struct {
	CosmosObj
	Name        string `gorm:"uniqueIndex"`
	Description string
	TeamID      *int
	Team        *Team `gorm:"foreignKey:TeamID"`
}
