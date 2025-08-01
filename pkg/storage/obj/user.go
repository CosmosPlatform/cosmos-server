package obj

type User struct {
	CosmosObj
	Email             string `gorm:"uniqueIndex"`
	Username          string `gorm:"unique"`
	EncryptedPassword string
	Role              string
	TeamID            *int
	Team              *Team `gorm:"foreignKey:TeamID"`
}
