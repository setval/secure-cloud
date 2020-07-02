package file

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"github.com/xxtea/xxtea-go/xxtea"
	"io/ioutil"
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
	b, err := ioutil.ReadAll(f.Body)
	if err != nil {
		return err
	}
	if err := f.Body.Close(); err != nil {
		return err
	}

	encrFile := xxtea.Encrypt(b, []byte(f.Key))

	f.Body = nil

	var file bytes.Buffer
	_, err = file.Write(encrFile)
	if err != nil {
		return err
	}
	f.Body = ioutil.NopCloser(bytes.NewReader(file.Bytes()))

	return nil
}

func (f *File) Decrypt() error {
	b, err := ioutil.ReadAll(f.Body)
	if err != nil {
		return err
	}
	if err := f.Body.Close(); err != nil {
		return err
	}

	decrFile := xxtea.Decrypt(b, []byte(f.Key))

	f.Body = nil
	var file bytes.Buffer
	_, err = file.Write(decrFile)
	if err != nil {
		return err
	}
	f.Body = ioutil.NopCloser(bytes.NewReader(file.Bytes()))

	return nil
}
