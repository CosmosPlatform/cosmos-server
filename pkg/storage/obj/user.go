package obj

type User struct {
	Email             string `db:"email"`
	Username          string `db:"username"`
	EncryptedPassword string `db:"encrypted_password"`
	Role              string `db:"role"`
}
