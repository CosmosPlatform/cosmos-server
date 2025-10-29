package model

type Token struct {
	Name           string
	EncryptedValue string
	Team           *Team
}

type TokenUpdate struct {
	Name  *string
	Value *string
}
