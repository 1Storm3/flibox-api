package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	appConfig "kbox-api/internal/config"
	"kbox-api/shared/httperror"
)

type S3Service struct {
	client    *s3.Client
	appConfig *appConfig.Config
}

func NewS3Service(appConfig *appConfig.Config) (*S3Service, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appConfig.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			appConfig.S3.AccessKey,
			appConfig.S3.SecretKey,
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
		o.BaseEndpoint = aws.String(appConfig.S3.Endpoint)
		o.UsePathStyle = true
	})

	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(appConfig.S3.Bucket),
	})
	if err != nil {
		return nil, httperror.New(
			http.StatusInternalServerError,
			"Не удалось подключиться к S3: "+err.Error(),
		)
	}

	return &S3Service{client: client, appConfig: appConfig}, nil
}

func (s *S3Service) UploadFile(ctx context.Context, key string, file []byte) (string, error) {
	if s.client == nil {
		return "", httperror.New(http.StatusInternalServerError, "S3 не инициализирован")
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.appConfig.S3.Bucket),
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

	url := fmt.Sprintf("%s/%s", s.appConfig.S3.Domain, key)
	return url, nil
}
