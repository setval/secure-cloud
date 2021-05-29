package file

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"github.com/xxtea/xxtea-go/xxtea"
	"golang.org/x/crypto/pbkdf2"
)

func GenKey(id uuid.UUID) (string, error) {
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := pbkdf2.Key([]byte(id.String()), salt, 10000, 2048, sha256.New)

	return base64.URLEncoding.EncodeToString(salt) + base64.URLEncoding.EncodeToString(hash), nil
}

func Encrypt(key string, b []byte) []byte {
	return xxtea.Encrypt(b, []byte(key))
}

func Decrypt(key string, b []byte) []byte {
	d := xxtea.Decrypt(b, []byte(key))
	return d
}
