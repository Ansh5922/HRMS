package storage

import (
"context"
"mime/multipart"

"github.com/aws/aws-sdk-go-v2/aws"
"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
Client *s3.Client
Bucket string
}

func (s *Storage) UploadFile(ctx context.Context, key string, file multipart.File, size int64, contentType string) (string, error) {
_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
Bucket:        aws.String(s.Bucket),
Key:           aws.String(key),
Body:          file,
ContentLength: aws.Int64(size),
ContentType:   aws.String(contentType),
})
if err != nil {
return "", err
}
return "https://" + s.Bucket + "/" + key, nil
}
