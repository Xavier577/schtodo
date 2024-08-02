package passwd

import "golang.org/x/crypto/bcrypt"

func Hash(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	return string(bytes), err
}

func Compare(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}
