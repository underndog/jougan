package aws_cloud

type AWSCloud interface {
	CreatePreSignedURL(bucket string, key string) (string, error)
	DownloadS3FileToMemory(bucket string, key string) ([]byte, int64, error)
}
