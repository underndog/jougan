package handler

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"jougan/helper"
	"jougan/helper/monitor"
	"jougan/log"
	"jougan/model"
	"net/http"
	"os"
	"path"
	"time"
)

type InspectDiskHandler struct {
	Monitoring monitor.Monitoring
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

	url := helper.GetEnvOrDefault("DOWNLOAD_URL", "https://www.dundeecity.gov.uk/sites/default/files/publications/civic_renewal_forms.zip")
	filePath := helper.GetEnvOrDefault("SAVE_TO_LOCATION", "save/dynamicSize.bin")

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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading the response body:", err)
		return
	}

	dataSize := len(data)
	id.Monitoring.FileSizeMonitor(path.Base(url), float64(dataSize))
	//fmt.Println("Size of the downloaded file:", helper.FormatSize(dataSize))
	//
	elapsedDownload := time.Since(startDownload).Seconds()
	downloadSpeed := float64(dataSize) / elapsedDownload
	//fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
	//fmt.Printf("Download speed: %f KB/s\n", downloadSpeed/1024)
	id.Monitoring.SpeedMonitor(path.Base(url), "download", downloadSpeed, elapsedDownload)

	//// Save
	startSave := time.Now()
	out, err := os.Create(filePath)
	if err != nil {
		log.Error("Error creating the file:", err)
		return
	}
	_, err = io.Copy(out, bytes.NewReader(data))
	out.Close()
	if err != nil {
		log.Error("Error saving the file:", err)
		return
	}

	elapsedSave := time.Since(startSave).Seconds()
	saveSpeed := float64(dataSize) / elapsedSave
	//fmt.Printf("Time taken to save the file: %f seconds\n", elapsedSave)
	//fmt.Printf("Save speed: %f KB/s\n", saveSpeed/1024)
	id.Monitoring.SpeedMonitor(path.Base(url), "save", saveSpeed, elapsedSave)

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
	id.Monitoring.SpeedMonitor(path.Base(url), "delete", deleteSpeed, elapsedDelete)
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
