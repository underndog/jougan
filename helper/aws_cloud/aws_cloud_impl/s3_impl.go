package aws_cloud_impl

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
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

func (ac *AWSConfiguration) DownloadS3FileToMemory(bucket string, key string) ([]byte, int64, error) {
	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return nil, 0, err
	}
	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Get the object from S3
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		log.Errorf("unable to download item from bucket %q, %v", bucket, err)
		return nil, 0, err
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, result.Body)
	if err != nil {
		log.Errorf("unable to read object into buffer, %v", err)
		return nil, 0, err
	}
	return buf.Bytes(), n, nil
}
