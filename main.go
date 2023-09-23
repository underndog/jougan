package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "https://speed.hetzner.de/100MB.bin"
	filePath := "save/100MB.bin"

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

	elapsedDownload := time.Since(startDownload).Seconds()
	downloadSpeed := float64(len(data)) / elapsedDownload
	fmt.Printf("Time taken to download the file: %f seconds\n", elapsedDownload)
	fmt.Printf("Download speed: %f MB/s\n", downloadSpeed/1e6)

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
	saveSpeed := float64(len(data)) / elapsedSave
	fmt.Printf("Time taken to save the file: %f seconds\n", elapsedSave)
	fmt.Printf("Save speed: %f MB/s\n", saveSpeed/1e6)

	// Delete
	startDelete := time.Now()
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting the file:", err)
		return
	}
	elapsedDelete := time.Since(startDelete).Seconds()
	deleteSpeed := float64(len(data)) / elapsedDelete
	fmt.Printf("Time taken to delete the file: %f seconds\n", elapsedDelete)
	fmt.Printf("Delete speed: %f MB/s\n", deleteSpeed/1e6)
}
