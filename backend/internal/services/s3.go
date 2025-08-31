package services

import (
    "bytes"
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
    bucketName = "vendor-renewals"
    s3Client   *s3.Client
)

func init() {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
    if err != nil {
        panic("failed to load AWS config: " + err.Error())
    }
    s3Client = s3.NewFromConfig(cfg)
}

func UploadToS3(ctx context.Context, key string, data []byte) error {
    _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key), // prefix
        Body:   bytes.NewReader(data),
    })
    return err
}
