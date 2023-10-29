package model

type DownloadFile struct {
	DownloadURL string `json:"downloadUrl,omitempty"` // tags
	SaveTo      string `json:"saveTo,omitempty"`
}
