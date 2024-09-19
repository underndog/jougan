package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321_+="

// Function to generate a random string of length n
func randomString(n int) string {
	// Slice to hold the generated runes
	result := make([]rune, n)
	// Convert randomStringSource to a slice of runes
	sourceRunes := []rune(randomStringSource)
	// Get the length of the source runes
	sourceLen := big.NewInt(int64(len(sourceRunes)))

	// Loop to generate each random character
	for i := range result {
		// Generate a random index within the bounds of sourceRunes
		randomIndex, err := rand.Int(rand.Reader, sourceLen)
		if err != nil {
			panic(err) // Handle error, it shouldn't happen often
		}
		// Assign the random rune to the result
		result[i] = sourceRunes[randomIndex.Int64()]
	}

	// Return the generated string
	return string(result)
}

// Generate a unique file path for each pod using a random value (e.g., timestamp or UUID)
func AppendRandomToFilename(filePath string) string {
	// Get current time and format it as hhmm-ddmmyyyy
	currentTime := time.Now().Format("150405-02-01-2006")
	// Generate a random string of length 8
	randomStr := randomString(8)

	// Get the directory and file name
	dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)

	// Split the file name into name and extension
	ext := filepath.Ext(baseName)             // Get the file extension
	name := strings.TrimSuffix(baseName, ext) // Remove the extension from the name

	// Append the timestamp to the file name
	newName := fmt.Sprintf("%s-%s-%s%s", name, currentTime, randomStr, ext)

	// Return the updated path
	return filepath.Join(dir, newName)
}
