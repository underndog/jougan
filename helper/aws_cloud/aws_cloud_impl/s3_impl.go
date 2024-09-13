package aws_cloud_impl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"jougan/log"
	"os"
	"sort"
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
		log.Debug("PART_SIZE_MB env: " + partSizeStr)
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

// UploadPartResult holds the result of a part upload
type UploadPartResult struct {
	Part types.CompletedPart
	Err  error
}

func (ac *AWSConfiguration) UploadFileToS3(filePath string, bucket string, key string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("failed to open file %q, %v", filePath, err)
		return err
	}
	defer file.Close()

	cfg, err := ac.accessAWSCloud()
	if err != nil {
		log.Error("Error creating new access key: %v", err)
		return err
	}
	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	partSize, err := getPartSizeFromEnv()
	if err != nil {
		return err
	}

	if partSize == 0 {
		// Single-part upload if PART_SIZE_MB is not set
		log.Debug("Uploading file as a single upload")

		// Create a new S3 uploader manager from the client
		uploader := manager.NewUploader(client)

		// Upload the file to the S3 bucket
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key), // S3 key name
			Body:   file,
		})
		if err != nil {
			log.Errorf("failed to upload file to S3, %v", err)
			return err
		}
		log.Debugf("File uploaded to S3 successfully: %s/%s", bucket, key)
		return nil
	}
	// Multipart upload if part size is specified
	log.Debugf("Multipart upload with part size %d bytes", partSize)

	// Step 1: Initiate Multipart Upload
	createOutput, err := client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Errorf("unable to initiate multipart upload: %v", err)
		return err
	}

	uploadID := *createOutput.UploadId
	log.Debugf("Multipart upload initiated, Upload ID: %s", uploadID)

	// Step 2: Set up the channels and WaitGroup
	var wg sync.WaitGroup
	partResults := make(chan UploadPartResult)
	buffer := make([]byte, partSize)
	partNum := int32(1)

	// Read and upload each part concurrently
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Errorf("error reading file: %v", err)
			return err
		}
		// Increment WaitGroup counter for each part being uploaded
		wg.Add(1)

		// Upload each part concurrently in a separate goroutine
		go func(partNumber int32, data []byte) {
			defer wg.Done()

			// Upload the part
			uploadPartOutput, err := client.UploadPart(context.TODO(), &s3.UploadPartInput{
				Bucket:     aws.String(bucket),
				Key:        aws.String(key),
				PartNumber: aws.Int32(partNumber),
				UploadId:   aws.String(uploadID),
				Body:       bytes.NewReader(data),
			})

			// Send the result back to the channel
			partResults <- UploadPartResult{
				Part: types.CompletedPart{
					ETag:       uploadPartOutput.ETag,
					PartNumber: aws.Int32(partNumber),
				},
				Err: err,
			}
		}(partNum, buffer[:bytesRead]) // pass partNum and part data as parameters to the goroutine

		partNum++
	}

	// Close the channel once all parts are done
	go func() {
		wg.Wait()
		close(partResults)
	}()

	// Collect the results from the channel and check for errors
	var completedParts []types.CompletedPart
	for result := range partResults {
		if result.Err != nil {
			log.Fatalf("Error uploading part: %v", result.Err)
			return result.Err
		}
		completedParts = append(completedParts, result.Part)
		log.Debugf("Uploaded part %d", *result.Part.PartNumber)
	}

	// Step 3: Sort parts by PartNumber to avoid InvalidPartOrder error
	sort.Slice(completedParts, func(i, j int) bool {
		return *completedParts[i].PartNumber < *completedParts[j].PartNumber
	})

	// Step 4: Complete Multipart Upload
	_, err = client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		log.Fatalf("Failed to complete multipart upload: %v", err)
		return err
	}

	log.Debugf("File uploaded successfully using multipart upload: %s/%s", bucket, key)
	return nil
}

// getPartSizeFromEnv gets the part size from an environment variable (in MB)
func getPartSizeFromEnv() (int64, error) {
	partSizeStr := os.Getenv("PART_SIZE_MB")
	if partSizeStr == "" {
		return 0, nil // No part size specified, treat as single file upload
	}

	// Convert the part size from MB to bytes
	partSizeMB, err := strconv.ParseInt(partSizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid PartSize: %v", err)
	}
	return partSizeMB * 1024 * 1024, nil // Convert MB to bytes
}
