package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"jougan/helper"
	"jougan/helper/aws_cloud"
	"jougan/helper/monitor"
	"jougan/log"
	"jougan/model"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type InspectDiskHandler struct {
	Monitoring monitor.Monitoring
	AWSCloud   aws_cloud.AWSCloud
}

// SignUp godoc
// @Summary Test Donwnload File
// @Description	Measure Download File and Save to Disk
// @Tags inspect
// @Accept  json
// @Produce  json
// @Param data body model.DownloadFile true "measure"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /inspect/download-url [post]
func (id *InspectDiskHandler) HandlerInspectDownloadFile(c echo.Context) error {
	req := model.DownloadFile{}
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	measure, err := downloadFile(req)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Successful Process",
		Data:       measure,
	})
}

func (id *InspectDiskHandler) DiskHandler() {
	log.Debug("Begin to measure the dowloading file - Debug")

	//os.Setenv("DOWNLOAD_FROM_S3_BUCKET", "ahihi-09262023")
	//os.Setenv("DOWNLOAD_FROM_S3_KEY", "4def6e02f687bd0ea3544a417ac2080b88d9a3b0511419da4ad7776167fc8545")

	downloadType := helper.GetEnvOrDefault("DOWNLOAD_TYPE", "")

	var data []byte
	var dataSize int
	var fileName string

	if downloadType == "AWS-S3-SDK" {
		log.Debug("Download file from S3 by SDK")
		s3Bucket, _ := os.LookupEnv("DOWNLOAD_FROM_S3_BUCKET")
		S3Key, _ := os.LookupEnv("DOWNLOAD_FROM_S3_KEY")

		startDownload := time.Now()
		fileData, _, err := id.AWSCloud.DownloadS3FileToMemory(s3Bucket, S3Key)
		if err != nil {
			log.Error("Error downloading the file:", err)
			return
		}
		data = fileData
		dataSize = len(data)
		// Extract the file name from the object key
		fileName = filepath.Base(S3Key)
		id.Monitoring.FileSizeMonitor(fileName, float64(dataSize))

		elapsedDownload := time.Since(startDownload).Seconds()
		downloadSpeed := float64(dataSize) / elapsedDownload
		//fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
		//fmt.Printf("Download speed: %f KB/s\n", downloadSpeed/1024)
		id.Monitoring.SpeedMonitor(fileName, "download", downloadSpeed, elapsedDownload)

	} else {
		// call get Function: getDownloadURL() in this file
		url, fetchedFileName := id.getDownloadURL()
		fileName = fetchedFileName // Set fileName from getDownloadURL

		// Download
		startDownload := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			log.Error("Error downloading the file:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("HTTP error:", resp.Status)
			return
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error reading the response body:", err)
			return
		}

		dataSize = len(data)
		id.Monitoring.FileSizeMonitor(fileName, float64(dataSize))
		//fmt.Println("Size of the downloaded file:", helper.FormatSize(dataSize))
		//
		elapsedDownload := time.Since(startDownload).Seconds()
		downloadSpeed := float64(dataSize) / elapsedDownload
		//fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
		//fmt.Printf("Download speed: %f KB/s\n", downloadSpeed/1024)
		id.Monitoring.SpeedMonitor(fileName, "download", downloadSpeed, elapsedDownload)
	}

	filePath := helper.GetEnvOrDefault("SAVE_TO_LOCATION", "save/dynamicSize.bin")

	//// Save
	startSave := time.Now()
	out, err := os.Create(filePath)
	if err != nil {
		log.Error("Error creating the file:", err)
		return
	}
	log.Debugf("Data size before saving: %d bytes", len(data))
	// Write data to the file
	_, err = io.Copy(out, bytes.NewReader(data))
	out.Close() // Done use defer out.Close() Because it will conflict delete file
	if err != nil {
		log.Error("Error saving the file:", err)
		return
	}

	elapsedSave := time.Since(startSave).Seconds()
	saveSpeed := float64(dataSize) / elapsedSave
	id.Monitoring.SpeedMonitor(fileName, "save", saveSpeed, elapsedSave)

	// Checksum File
	// Only calculate checksum if the environment variable is set
	sha256FromEnv := helper.GetEnvOrDefault("SHA-256-CHECKSUM", "")
	if sha256FromEnv != "" {
		// Calculate the SHA-256 checksum
		sha256Hasher := sha256.New()
		_, err := io.Copy(sha256Hasher, bytes.NewReader(data))
		if err != nil {
			log.Error("Error calculating checksum:", err)
			return
		}

		// Convert the SHA-256 binary to a hexadecimal string
		calculatedSHA256 := sha256Hasher.Sum(nil)
		calculatedSHA256String := hex.EncodeToString(calculatedSHA256)

		// Convert to Base64
		calculatedChecksumBase64 := base64.StdEncoding.EncodeToString([]byte(calculatedSHA256String))

		log.Debugf("Calculated SHA-256 (Base64): %s", calculatedChecksumBase64)

		// Encode the provided checksum from the environment variable to Base64
		Sha256Base64FromEnv := base64.StdEncoding.EncodeToString([]byte(sha256FromEnv))

		// Compare the calculated checksum with the checksum from the environment variable
		if calculatedChecksumBase64 == Sha256Base64FromEnv {
			log.Debug("The file is verified. SHA-256 matches the S3 checksum.")
		} else {
			log.Debugf("SHA-256 mismatch! S3 checksum: %s, Calculated SHA-256: %s", Sha256Base64FromEnv, calculatedChecksumBase64)
		}
	}

	// UploadFile
	uploadFileToS3, err := strconv.ParseBool(helper.GetEnvOrDefault("UPLOAD_FILE_TO_S3", "true"))
	if err != nil {
		log.Error(err)
	}
	if uploadFileToS3 == true {
		s3Bucket, _ := os.LookupEnv("DOWNLOAD_FROM_S3_BUCKET")
		S3Key, _ := os.LookupEnv("DOWNLOAD_FROM_S3_KEY")
		startUpload := time.Now()
		new_S3key := helper.AppendTimestampToFile(S3Key)
		err = id.AWSCloud.UploadFileToS3(filePath, s3Bucket, new_S3key)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Error(err)
			return
		}
		elapsedUpload := time.Since(startUpload).Seconds()
		uploadSpeed := float64(dataSize) / elapsedUpload
		id.Monitoring.SpeedMonitor(fileName, "upload", uploadSpeed, elapsedUpload)
	}

	// Delete
	startDelete := time.Now()
	err = os.Remove(filePath)
	if err != nil {
		log.Error("Error deleting the file:", err)
		return
	}
	elapsedDelete := time.Since(startDelete).Seconds()
	deleteSpeed := float64(dataSize) / elapsedDelete
	//fmt.Printf("Time taken to delete the file: %f seconds\n", elapsedDelete)
	//fmt.Printf("Delete speed: %f KB/s\n", deleteSpeed/1024)
	id.Monitoring.SpeedMonitor(fileName, "delete", deleteSpeed, elapsedDelete)
	log.Debug("Completed delete file: ", filePath)
}

func downloadFile(download model.DownloadFile) (model.MetricResponse, error) {

	measureResult := model.MetricResponse{}
	// Download
	startDownload := time.Now()
	resp, err := http.Get(download.DownloadURL)
	if err != nil {
		log.Error("Error downloading the file:", err)
		return measureResult, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("HTTP error:", resp.Status)
		return measureResult, fmt.Errorf("HTTP error: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading the response body:", err)
		return measureResult, err
	}

	dataSize := len(data)
	measureResult.FileSize = dataSize

	//fmt.Println("Size of the downloaded file:", helper.FormatSize(dataSize))
	//
	elapsedDownload := time.Since(startDownload).Seconds()
	downloadSpeed := float64(dataSize) / elapsedDownload
	//fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
	//fmt.Printf("Download speed: %f KB/s\n", downloadSpeed/1024)

	measureResult.DownloadTime = elapsedDownload
	measureResult.DownloadSpeed = downloadSpeed

	//// Save
	startSave := time.Now()
	out, err := os.Create(download.SaveTo)
	if err != nil {
		log.Error("Error creating the file:", err)
		return measureResult, err
	}
	_, err = io.Copy(out, bytes.NewReader(data))
	out.Close()
	if err != nil {
		log.Error("Error saving the file:", err)
		return measureResult, err
	}

	elapsedSave := time.Since(startSave).Seconds()
	saveSpeed := float64(dataSize) / elapsedSave
	//fmt.Printf("Time taken to save the file: %f seconds\n", elapsedSave)
	//fmt.Printf("Save speed: %f KB/s\n", saveSpeed/1024)
	measureResult.SaveTime = elapsedSave
	measureResult.SaveSpeed = saveSpeed

	// Delete
	startDelete := time.Now()
	err = os.Remove(download.SaveTo)
	if err != nil {
		log.Error("Error deleting the file:", err)
		return measureResult, err
	}
	elapsedDelete := time.Since(startDelete).Seconds()
	deleteSpeed := float64(dataSize) / elapsedDelete
	measureResult.DeleteTime = elapsedDelete
	measureResult.DeleteSpeed = deleteSpeed

	return measureResult, nil
}

/*
**
This function return Download URL and File Name
*/
func (id *InspectDiskHandler) getDownloadURL() (string, string) {
	s3Bucket, okBucket := os.LookupEnv("DOWNLOAD_FROM_S3_BUCKET")
	S3Key, okKey := os.LookupEnv("DOWNLOAD_FROM_S3_KEY")

	if !okBucket || !okKey {
		url := helper.GetEnvOrDefault("DOWNLOAD_URL", "https://www.dundeecity.gov.uk/sites/default/files/publications/civic_renewal_forms.zip")
		return url, path.Base(url)
	}

	log.Debug("Generate Pre-Signed URL of ", s3Bucket, " Bucket")
	presignedURL, err := id.AWSCloud.CreatePreSignedURL(s3Bucket, S3Key)
	if err != nil {
		log.Error(err)
		return "", ""
	}
	log.Debug(presignedURL)
	return presignedURL, S3Key

}
