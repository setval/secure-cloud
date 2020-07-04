package file

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"github.com/xxtea/xxtea-go/xxtea"
	"golang.org/x/crypto/pbkdf2"
)

func (f *File) genKey() error {
	f.ID = uuid.NewV4()

	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	hash := pbkdf2.Key([]byte(f.ID.String()), salt, 10000, 64, sha256.New)

	_, err = rand.Read(salt)
	if err != nil {
		return err
	}

	f.Key = base64.URLEncoding.EncodeToString(salt) + base64.URLEncoding.EncodeToString(hash)
	return nil
}

func (f *File) Encrypt() error {
	if err := f.genKey(); err != nil {
		return err
	}

	encrFile := xxtea.Encrypt(f.Bytes(), []byte(f.Key))
	f.Reset()

	if _, err := f.Write(encrFile); err != nil {
		return err
	}

	return nil
}

func (f *File) Decrypt() error {
	decrFile := xxtea.Decrypt(f.Bytes(), []byte(f.Key))
	f.Reset()

	_, err := f.Write(decrFile)
	return err
}
