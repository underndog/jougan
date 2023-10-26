package monitor

type Monitoring interface {
	FileSizeMonitor(filename string, fileSize float64)
	SpeedMonitor(filename string, activity string, speed float64, time float64)
}
