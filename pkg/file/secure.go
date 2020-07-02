package file

import (
	"crypto/sha1"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"github.com/xxtea/xxtea-go/xxtea"
)

func (f *File) genKey() error {
	f.ID = uuid.NewV4()
	hash := sha1.New()
	_, err := hash.Write([]byte(f.ID.String()))
	if err != nil {
		return err
	}
	f.Key = base64.URLEncoding.EncodeToString(hash.Sum(nil))
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
