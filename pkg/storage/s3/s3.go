//+build ignore

// Support for uploading and downloading files from S3 storage.

package s3

import (
	"github.com/DiscoreMe/SecureCloud/pkg/file"
	"github.com/DiscoreMe/SecureCloud/pkg/storage"
	"github.com/minio/minio-go/v6"
	"path"
)

type S3 struct {
	endpoint string
	bucket   string
	client   *minio.Client
}

func New(endpoint, bucket, accessKey, secretKey, location string) (*S3, error) {
	s3 := &S3{
		endpoint: endpoint,
		bucket:   bucket,
	}

	client, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		return nil, err
	}

	if ok, err := client.BucketExists(bucket); err != nil {
		return nil, err
	} else if !ok {
		if err := client.MakeBucket(bucket, location); err != nil {
			return nil, err
		}
	}

	s3.client = client

	return s3, nil
}

func (s *S3) Upload(f *file.File) error {
	_, err := s.client.PutObject(s.bucket, path.Join(storage.Folder, f.ID.String()), f, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3) Download(f *file.File) error {
	object, err := s.client.GetObject(s.bucket, path.Join(storage.Folder, f.ID.String()), minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()
	if _, err := f.WriteFromReader(object); err != nil {
		return err
	}
	return nil
}
