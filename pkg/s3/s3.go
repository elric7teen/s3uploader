package s3

import (
	"os"

	pkgAws "ghozi.com/prototype/s3uploader/pkg/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type repo struct {
	s3Client *s3.S3
	uploader S3Uploader
}

type S3Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func NewRepo(awsClient *pkgAws.Client) *repo {
	s3Client := s3.New(awsClient.Sess, awsClient.Config)
	uploader := s3manager.NewUploaderWithClient(s3Client)

	return &repo{
		s3Client: s3Client,
		uploader: uploader,
	}
}

// Upload : upload file into s3
func (r *repo) Upload(filepath, bucket, key string) error {

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}
