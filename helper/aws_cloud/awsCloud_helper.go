package aws_cloud

type AWSCloud interface {
	CreatePreSignedURL(bucket string, key string) (string, error)
}
