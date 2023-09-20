package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strings"
)

type Client struct {
	s3Config *aws.Config
	bucket   string
}

func NewS3Client(id, secret, token, endpoint, region, bucket string, disableSSL bool) *Client {
	s3Client := &Client{
		s3Config: &aws.Config{
			Credentials:      credentials.NewStaticCredentials(id, secret, token),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			DisableSSL:       aws.Bool(disableSSL),
			S3ForcePathStyle: aws.Bool(true),
		},
		bucket: bucket,
	}
	return s3Client
}

func (s *Client) PutObject(key, ext string, file io.ReadSeeker) (err error) {
	newSession, err := session.NewSession(s.s3Config)
	if err != nil {
		return fmt.Errorf("failed to create session, %s", err.Error())
	}
	contentType := fmt.Sprintf("image/%s", strings.TrimPrefix(ext, "."))
	_, err = s3.New(newSession).PutObject(&s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to put object, %s", err.Error())
	}
	return nil
}
