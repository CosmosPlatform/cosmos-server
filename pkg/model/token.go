package model

type Token struct {
	Name           string
	EncryptedValue string
	Team           *Team
}
