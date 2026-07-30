package privilege

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/phoenix-agent-go/internal/model"
)

// HashPassword generates an MD5 hash of "phoenix" + plaintext.
func HashPassword(password string) (string, error) {
	h := md5.Sum([]byte("phoenix" + password))
	return hex.EncodeToString(h[:]), nil
}

// CheckPassword compares a plaintext password against a stored MD5 hash.
func CheckPassword(hashed, password string) bool {
	h := md5.Sum([]byte("phoenix" + password))
	return hashed == hex.EncodeToString(h[:])
}

// GenerateRandomPassword returns an 8-character random password composed of
// uppercase letters, lowercase letters, and digits.
func GenerateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// IsValidUser checks whether a user is active.
// Returns true when user != nil AND status != 1 (status == 1 means DISABLED).
func IsValidUser(user *model.PrivilegeUser) bool {
	return user != nil && user.Status != 1
}
