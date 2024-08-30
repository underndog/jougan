package aws_cloud_impl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"jougan/log"
	"os"
	"strconv"
	"sync"
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

	// Get file size
	headObjectOutput, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, 0, err
	}
	fileSize := *headObjectOutput.ContentLength

	// Check for PART_SIZE environment variable
	partSizeStr := os.Getenv("PART_SIZE_MB")
	var partSize int64
	if partSizeStr != "" {
		partSizeMB, err := strconv.ParseInt(partSizeStr, 10, 64)
		if err != nil {
			log.Errorf("Invalid PART_SIZE value: %v", err)
			return nil, 0, err
		}
		partSize = partSizeMB * 1024 * 1024 // Convert MB to bytes
	}

	// If partSize is not set, download the entire file in one go
	if partSize <= 0 {
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

	// Calculate the number of parts
	numParts := (fileSize + partSize - 1) / partSize

	// Create a buffer to store the file content
	buf := make([]byte, fileSize)

	// Download parts concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, numParts)

	for i := int64(0); i < numParts; i++ {
		start := i * partSize
		end := start + partSize - 1
		if end > fileSize-1 {
			end = fileSize - 1
		}

		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()

			getObjectInput := &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
			}

			result, err := client.GetObject(context.TODO(), getObjectInput)
			if err != nil {
				log.Errorf("unable to download part from bucket %q, %v", bucket, err)
				errCh <- err
				return
			}
			defer result.Body.Close()

			// Read the part into the buffer at the correct offset
			partBuffer := make([]byte, end-start+1)
			_, err = io.ReadFull(result.Body, partBuffer)
			if err != nil {
				log.Errorf("unable to read object part into buffer, %v", err)
				errCh <- err
				return
			}

			// Copy the part into the correct position in the buffer
			copy(buf[start:], partBuffer)
		}(start, end)
	}

	// Wait for all parts to be downloaded
	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, 0, <-errCh // return the first error encountered
	}

	return buf, fileSize, nil
}
