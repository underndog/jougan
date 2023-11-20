package aws_cloud_impl

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"jougan/log"
	"time"
)

func (ac *AWSConfiguration) CreatePreSignedURL(bucket string, key string) (string, error) {
	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return "", err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	//// Define the object you want to create a pre-signed URL for
	//bucket := "your-bucket-name"
	//key := "your-object-key"

	// Define the duration that the pre-signed URL should be valid for
	//  example, 1 minutes
	expireTime := 1 * time.Minute

	//// Create a GetObject request
	//getObjectInput := &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(key),
	//}

	// Create a GetObject presigned URL
	presigner := s3.NewPresignClient(client)
	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		log.Fatalf("failed to presign request, %v", err)
	}

	// Output the presigned URL
	//fmt.Println("Presigned URL:", req.URL)
	return req.URL, nil
}
