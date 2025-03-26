package wasabi

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Conf struct {
	bucketName string
	region     string
	endpoint   string
	accessKey  string
	secretKey  string
}

type S3 interface {
	PutFile(localPath, remotePath string) error
	GetFile(remotePath, localPath string) error
	DelFile(remotePath string) error
	GetBucketName() string
	GetRegion() string
	GetEndpoint() string
	GetAccessKey() string
	GetSecretKey() string
}

func (c *s3Conf) GetBucketName() string { return c.bucketName }
func (c *s3Conf) GetRegion() string     { return c.region }
func (c *s3Conf) GetEndpoint() string   { return c.endpoint }
func (c *s3Conf) GetAccessKey() string  { return c.accessKey }
func (c *s3Conf) GetSecretKey() string  { return c.secretKey }

func (c *s3Conf) handler() *s3.S3 {
	config := &aws.Config{
		Region:      aws.String(c.GetRegion()),
		Endpoint:    aws.String(c.GetEndpoint()),
		Credentials: credentials.NewStaticCredentials(c.GetAccessKey(), c.GetSecretKey(), ""),
	}

	sess := session.Must(session.NewSession(config))
	return s3.New(sess)
}

func (c *s3Conf) PutFile(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	svc := c.handler()

	_, err = svc.PutObject(
		&s3.PutObjectInput{
			Bucket: aws.String(c.GetBucketName()),
			Key:    aws.String(remotePath),
			Body:   file,
		},
	)
	return err
}

func (c *s3Conf) GetFile(remotePath, localPath string) error {
	svc := c.handler()

	resp, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(c.GetBucketName()),
			Key:    aws.String(remotePath),
		},
	)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *s3Conf) DelFile(remotePath string) error {
	svc := c.handler()

	_, err := svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(c.GetBucketName()),
			Key:    aws.String(remotePath),
		},
	)
	return err
}

func (c *s3Conf) GetAllBuckets() ([]string, error) {
	svc := c.handler()

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	var buckets []string
	for _, bucket := range result.Buckets {
		buckets = append(buckets, aws.StringValue(bucket.Name))
	}

	return buckets, nil
}

func (c *s3Conf) GetAllFilesFromBucket(bucketName string) ([]string, error) {
	if bucketName == "" {
		bucketName = c.GetBucketName()
	}

	svc := c.handler()

	result, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range result.Contents {
		files = append(files, aws.StringValue(obj.Key))
	}

	return files, nil
}
