package obj

type User struct {
	Email             string `bson:"email"`
	Username          string `bson:"username"`
	EncryptedPassword string `bson:"encryptedPassword"`
	Role              string `bson:"role"`
}
