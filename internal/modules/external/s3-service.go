package external

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	appConfig "github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var _ S3ServiceInterface = (*S3Service)(nil)

type S3ServiceInterface interface {
	UploadFile(ctx context.Context, key string, file []byte) (string, error)
	DeleteFile(ctx context.Context, key string) error
}

type S3Service struct {
	client *s3.Client
	appCfg *appConfig.Config
}

func NewS3Service(appCfg *appConfig.Config) (*S3Service, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appCfg.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			appCfg.S3.AccessKey,
			appCfg.S3.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			"Не удалось загрузить конфигурацию для AWS SDK",
		)
	}

	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(appCfg.S3.Endpoint)
		o.UsePathStyle = true
	})

	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(appCfg.S3.Bucket),
	})
	if err != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			"Не удалось подключиться к S3: "+err.Error(),
		)
	}

	return &S3Service{client: client, appCfg: appCfg}, nil
}

func (s *S3Service) UploadFile(ctx context.Context, key string, file []byte) (string, error) {
	if s.client == nil {
		return "", httperror.New(http.StatusInternalServerError, "S3 не инициализирован")
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.appCfg.S3.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpeg"),
		ACL:         types.ObjectCannedACLPublicReadWrite,
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			"Не удалось загрузить файл: "+err.Error(),
		)
	}

	url := fmt.Sprintf("%s/%s", s.appCfg.S3.Domain, key)
	return url, nil
}

func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	if s.client == nil {
		return httperror.New(http.StatusInternalServerError, "S3 не инициализирован")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.appCfg.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			"Не удалось удалить файл: "+err.Error(),
		)
	}
	return nil
}
