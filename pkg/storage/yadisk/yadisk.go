package yadisk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DiscoreMe/SecureCloud/config"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/nikitaksv/yandex-disk-sdk-go"
	"io"
	"io/ioutil"
	"net/http"
)

type YaDisk struct{
	client yadisk.YaDisk
}


func New() *YaDisk {
	cfg := config.NewConfig()
	client, _ := yadisk.NewYaDisk(context.Background(),http.DefaultClient, &yadisk.Token{AccessToken: cfg.YadiskToken})
	return &YaDisk{client: client}
}

func (y *YaDisk) performUpload(f *file.File) error {
	link, err := y.client.GetResourceUploadLink(f.ID.String(), nil, true)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(f.FileBuffer.Bytes())
	_, err = y.client.PerformUpload(link, buf)
	return err
}

func (y *YaDisk) Upload(f *file.File) error {
	_, err := f.Write(f.Bytes())
	if err != nil {
		return err
	}
	err = y.performUpload(f)
	return err
}

func (y *YaDisk) performDownload(f *file.File) ([]byte, error){
	filename := f.ID.String()
	link, err := y.client.GetResourceDownloadLink(filename, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(link.Href)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	downloaded, err := ioutil.ReadAll(resp.Body)
	return downloaded, err
}

func (y *YaDisk) Download(f *file.File) error {
	body, err := y.performDownload(f)
	if err != nil {
		return err
	}
	_, err = f.FileBuffer.Write(body)
	return err
}