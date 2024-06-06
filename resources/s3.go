package resources

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"passVault/dtos"
)

type S3 struct {
	service *s3.S3
}

func NewS3() S3 {
	sess := session.Must(session.NewSession())

	return S3{
		service: s3.New(sess),
	}
}

func (s S3) Push(ctx context.Context, bucket, path string, data io.ReadSeeker) error {
	_, err := s.service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   data,
		Key:    &path,
		Bucket: &bucket,
	})
	return err
}

// List method will list down all the objects in the provided bucket and path
func (s S3) List(ctx context.Context, bucket, path string) ([]dtos.S3Object, error) {
	var objects []dtos.S3Object
	err := s.service.ListObjectsV2PagesWithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &path,
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range page.Contents {
			objects = append(objects, dtos.S3Object{
				Key:          *object.Key,
				Size:         *object.Size,
				LastModified: *object.LastModified,
			})
		}
		return !lastPage
	})
	return objects, err
}

// Delete method will Delete the object from the provided bucket and path
func (s S3) Delete(ctx context.Context, bucket, path string) error {
	_, err := s.service.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})
	return err
}

// DeleteBulk method will Delete the objects from the provided bucket and paths
func (s S3) DeleteBulk(ctx context.Context, bucket string, paths []string) error {
	var objects []*s3.ObjectIdentifier
	for _, path := range paths {
		objects = append(objects, &s3.ObjectIdentifier{
			Key: &path,
		})
	}
	_, err := s.service.DeleteObjectsWithContext(ctx, &s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: &s3.Delete{
			Objects: objects,
		},
	})
	return err
}
