package wasabi

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (c *s3Conf) handler() *s3.Client {
	cfg := aws.Config{
		Region:      c.GetRegion(),
		Credentials: credentials.NewStaticCredentialsProvider(c.GetAccessKey(), c.GetSecretKey(), ""),
	}

	if c.GetEndpoint() != "" {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               c.GetEndpoint(),
					SigningRegion:     c.GetRegion(),
					HostnameImmutable: true,
				}, nil
			})
	}

	return s3.NewFromConfig(cfg)
}

func (c *s3Conf) PutFile(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	svc := c.handler()

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.GetBucketName()),
		Key:    aws.String(remotePath),
		Body:   file,
	})
	return err
}

func (c *s3Conf) GetFile(remotePath, localPath string) error {
	svc := c.handler()

	resp, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.GetBucketName()),
		Key:    aws.String(remotePath),
	})
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

	_, err := svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(c.GetBucketName()),
		Key:    aws.String(remotePath),
	})
	return err
}

func (c *s3Conf) GetAllBuckets() ([]string, error) {
	svc := c.handler()

	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets []string
	for _, bucket := range result.Buckets {
		if bucket.Name != nil {
			buckets = append(buckets, *bucket.Name)
		}
	}

	return buckets, nil
}

func (c *s3Conf) GetAllFilesFromBucket(bucketName string) ([]string, error) {
	if bucketName == "" {
		bucketName = c.GetBucketName()
	}

	svc := c.handler()

	result, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range result.Contents {
		if obj.Key != nil {
			files = append(files, *obj.Key)
		}
	}

	return files, nil
}
