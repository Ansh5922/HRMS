package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func EmpCodeFromName(firstName, lastName string, seq int) string {
	f := strings.ToUpper(string([]rune(firstName)[0]))
	l := strings.ToUpper(string([]rune(lastName)[0]))
	return "EMP" + f + l + fmt.Sprintf("%04d", seq)
}
