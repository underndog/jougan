package helper

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func FormatSize(bytes int) string {
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

// Function to modify file name by appending timestamp in the format hhmm-ddmmyyyy
func AppendTimestampToFile(filePath string) string {
	// Get current time and format it as hhmm-ddmmyyyy
	currentTime := time.Now().Format("1504-02012006")

	// Get the directory and file name
	dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)

	// Split the file name into name and extension
	ext := filepath.Ext(baseName)             // Get the file extension
	name := strings.TrimSuffix(baseName, ext) // Remove the extension from the name

	// Append the timestamp to the file name
	newName := fmt.Sprintf("%s-%s%s", name, currentTime, ext)

	// Return the updated path
	return filepath.Join(dir, newName)
}

// Generate a unique file path for each pod using a random value (e.g., timestamp or UUID)
func AppendRandomToFilename(filePath string) string {
	// Get current time and format it as hhmm-ddmmyyyy
	randomValue := time.Now().UnixNano() // Or use a UUID generator if available

	// Get the directory and file name
	dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)

	// Split the file name into name and extension
	ext := filepath.Ext(baseName)             // Get the file extension
	name := strings.TrimSuffix(baseName, ext) // Remove the extension from the name

	// Append the timestamp to the file name
	newName := fmt.Sprintf("%s-%d%s", name, randomValue, ext)

	// Return the updated path
	return filepath.Join(dir, newName)
}
