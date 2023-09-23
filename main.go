package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func formatSize(bytes int) string {
	const (
		B  = 1.0
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes < int(KB):
		return fmt.Sprintf("%d Bytes", bytes)
	case bytes < int(MB):
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	case bytes < int(GB):
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	default:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	}
}

func main() {
	url := "https://www.dundeecity.gov.uk/sites/default/files/publications/civic_renewal_forms.zip"
	filePath := "save/dynamicSize.bin"

	// Download
	startDownload := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading the file:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP error:", resp.Status)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return
	}

	dataSize := len(data)
	fmt.Println("Size of the downloaded file:", formatSize(dataSize))

	elapsedDownload := time.Since(startDownload).Seconds()
	downloadSpeed := float64(dataSize) / elapsedDownload
	fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
	fmt.Printf("Download speed: %f KB/s\n", downloadSpeed/1024)

	// Save
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
	fmt.Printf("Time taken to save the file: %f seconds\n", elapsedSave)
	fmt.Printf("Save speed: %f KB/s\n", saveSpeed/1024)

	// Delete
	startDelete := time.Now()
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting the file:", err)
		return
	}
	elapsedDelete := time.Since(startDelete).Seconds()
	deleteSpeed := float64(dataSize) / elapsedDelete
	fmt.Printf("Time taken to delete the file: %f seconds\n", elapsedDelete)
	fmt.Printf("Delete speed: %f KB/s\n", deleteSpeed/1024)
}
