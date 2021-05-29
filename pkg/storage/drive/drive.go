//+build ignore

// Support for uploading and downloading files from Google drive.

package drive

import (
	"context"
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/DiscoreMe/SecureCloud/pkg/storage"
	"golang.org/x/oauth2"
	gdrive "google.golang.org/api/drive/v3"
	"net/http"
)

const (
	credentialsPath = "credentials.json"
	tokenPath       = "token.json"
)

type Drive struct {
	client   *http.Client
	svr      *gdrive.Service
	parentID string
}

func New() (*Drive, error) {
	d := &Drive{}

	if err := d.init(); err != nil {
		return nil, err
	}

	filesListResp, err := d.svr.Files.List().Do()
	if err != nil {
		return nil, err
	}
	for _, gfile := range filesListResp.Files {
		if gfile.Name == storage.Folder {
			d.parentID = gfile.Id
		}
	}

	if d.parentID == "" {
		fg := &gdrive.File{
			Name:     storage.Folder,
			MimeType: "application/vnd.google-apps.folder",
			Parents:  []string{"root"},
		}

		createResp, err := d.svr.Files.Create(fg).Do()
		if err != nil {
			return nil, err
		}
		d.parentID = createResp.Id
	}

	return d, nil
}

func (d *Drive) Upload(f *file.File) error {

	fg := &gdrive.File{
		Name:    f.ID.String(),
		Parents: []string{d.parentID},
	}
	_, err := d.svr.Files.Create(fg).Media(f).Do()
	return err
}

func (d *Drive) Download(f *file.File) error {
	filename := f.ID.String()
	filesListResp, err := d.svr.Files.List().Do()
	if err != nil {
		return err
	}
	for _, gfile := range filesListResp.Files {
		if gfile.Name == filename {
			httpResp, err := d.svr.Files.Get(gfile.Id).Download()
			if err != nil {
				return err
			}
			if _, err := f.WriteFromReader(httpResp.Body); err != nil {
				httpResp.Body.Close()
				return err
			}
			httpResp.Body.Close()
			break
		}
	}
	return nil
}

func (d *Drive) init() error {
	srvConfig, err := parseConfig(credentialsPath)
	if err != nil {
		return err
	}

	client, err := initClient(tokenPath, srvConfig)
	if err != nil {
		return err
	}
	d.client = client

	svr, err := gdrive.New(d.client)
	if err != nil {
		return err
	}
	d.svr = svr

	return nil
}

func initClient(tokenPath string, config *oauth2.Config) (*http.Client, error) {
	tokFile := tokenPath
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background(), tok)
	return client, nil
}
