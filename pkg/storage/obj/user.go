package obj

type User struct {
	ID                string `bson:"_id,omitempty"`
	Email             string `bson:"email"`
	Username          string `bson:"username"`
	EncryptedPassword string `bson:"encryptedPassword"`
	Role              string `bson:"role"`
}
