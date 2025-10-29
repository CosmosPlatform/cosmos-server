package obj

type Token struct {
	CosmosObj
	Name           string
	EncryptedValue string
	TeamID         int
	Team           *Team `gorm:"foreignKey:TeamID"`
}
