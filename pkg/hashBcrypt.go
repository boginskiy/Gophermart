package pkg

import "golang.org/x/crypto/bcrypt"

func GenerateHash(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashByte), err

}

func CompareHashAndPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
