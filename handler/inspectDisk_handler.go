package handler

import (
	"bytes"
	"fmt"
	"io"
	"jougan/helper/monitor"
	"jougan/log"
	"net/http"
	"os"
	"path"
	"time"
)

type InspectDiskHandler struct {
	Monitoring monitor.Monitoring
}

func (id *InspectDiskHandler) DiskHandler() {

	log.Info("Begin to measure the dowloading file")

	url := "https://www.dundeecity.gov.uk/sites/default/files/publications/civic_renewal_forms.zip"
	filePath := "save/dynamicSize.bin"

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
		fmt.Println("Error creating the file:", err)
		return
	}
	_, err = io.Copy(out, bytes.NewReader(data))
	out.Close()
	if err != nil {
		fmt.Println("Error saving the file:", err)
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
		fmt.Println("Error deleting the file:", err)
		return
	}
	elapsedDelete := time.Since(startDelete).Seconds()
	deleteSpeed := float64(dataSize) / elapsedDelete
	//fmt.Printf("Time taken to delete the file: %f seconds\n", elapsedDelete)
	//fmt.Printf("Delete speed: %f KB/s\n", deleteSpeed/1024)
	id.Monitoring.SpeedMonitor(path.Base(url), "delete", deleteSpeed, elapsedDelete)
}
