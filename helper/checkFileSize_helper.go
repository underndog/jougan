package helper

import "fmt"

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
